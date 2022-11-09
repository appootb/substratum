package configure

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/appootb/substratum/v2/configure"
	ictx "github.com/appootb/substratum/v2/internal/context"
	"github.com/appootb/substratum/v2/logger"
)

func Init() {
	if configure.CallbackImplementor() == nil {
		configure.RegisterCallbackImplementor(newCallback())
	}
	if configure.BackendImplementor() == nil {
		configure.RegisterBackendImplementor(newDebug())
	}
	if configure.Implementor() == nil {
		configure.RegisterImplementor(&Configure{})
	}
}

const (
	TagDefault = "default"
	TagComment = "comment"

	ConfigPrefix = "config"
)

var (
	staticType  = reflect.TypeOf((*configure.StaticType)(nil)).Elem()
	dynamicType = reflect.TypeOf((*configure.DynamicType)(nil)).Elem()
)

type ConfigItem struct {
	Type    string `json:"type"`
	Schema  string `json:"schema"`
	Value   string `json:"value"`
	Comment string `json:"comment"`
}

func (c ConfigItem) String() string {
	v, _ := json.Marshal(c)
	return string(v)
}

type ConfigItems map[string]*ConfigItem

func (c ConfigItems) Add(items ConfigItems) {
	for k, v := range items {
		c[k] = v
	}
}

func (c ConfigItems) KVs(path string) map[string]string {
	kvs := make(map[string]string, len(c))
	for key, item := range c {
		kvs[path+key] = item.String()
	}
	return kvs
}

type Configure struct{}

// Register the configuration pointer.
func (m *Configure) Register(component string, v interface{}, opts ...configure.Option) error {
	options := configure.EmptyOptions()
	for _, o := range opts {
		o(options)
	}
	//
	cfg := reflect.ValueOf(v)
	if cfg.Kind() != reflect.Ptr || cfg.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(cfg)}
	}

	basePath := fmt.Sprintf("%s/%s/%s/", options.Path, ConfigPrefix, component)
	// Create/update default config value if not exist.
	if options.AutoCreation {
		err := m.migrate(basePath, reflect.TypeOf(v))
		if err != nil {
			return err
		}
	}

	// Initialize the config.
	version, err := m.getConfig(basePath, m.configElem(cfg), false)
	if err != nil {
		return err
	}
	// Watch for config updated.
	evtChan, err := configure.BackendImplementor().Watch(basePath, version, true)
	if err != nil {
		return err
	}

	go m.watchEvent(basePath, evtChan, m.configElem(cfg))
	return nil
}

func (m *Configure) configElem(v reflect.Value) reflect.Value {
	if v.Type().Kind() == reflect.Ptr {
		return m.configElem(v.Elem())
	}
	return v
}

func (m *Configure) migrate(basePath string, t reflect.Type) error {
	provider := configure.BackendImplementor()
	pairs, err := provider.Get(basePath, true)
	if err != nil {
		return err
	}
	reflectKVs := m.parseConfig(t, "").KVs(basePath)
	for _, pair := range pairs.KVs {
		delete(reflectKVs, pair.Key)
	}
	for k, v := range reflectKVs {
		if err = provider.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (m *Configure) parseConfig(t reflect.Type, baseName string) ConfigItems {
	if t.Kind() == reflect.Ptr {
		return m.parseConfig(t.Elem(), baseName)
	}
	return m.parseConfigItems(t, baseName)
}

func (m *Configure) parseConfigItems(t reflect.Type, baseName string) ConfigItems {
	items := ConfigItems{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if m.isSupportedType(field.Type, 0) {
			items[baseName+field.Name] = &ConfigItem{
				Type:    strings.ReplaceAll(field.Type.String(), "*", ""),
				Schema:  "", // TODO
				Value:   m.formatDefaultValue(field.Type, field.Tag),
				Comment: field.Tag.Get(TagComment),
			}
		} else if field.Type.Kind() == reflect.Ptr {
			items.Add(m.parseConfigItems(field.Type, baseName))
		} else if field.Type.Kind() == reflect.Struct {
			items.Add(m.parseConfigItems(field.Type, baseName+field.Name+"/"))
		} else {
			panic("substratum: unsupported field type:" + field.Type.String())
		}
	}
	return items
}

func (m *Configure) formatDefaultValue(t reflect.Type, tag reflect.StructTag) string {
	val := tag.Get(TagDefault)
	if m.isSliceOrMap(t) && !strings.Contains(val, ";") {
		val = strings.ReplaceAll(val, ",", ";")
	}
	return val
}

func (m *Configure) isSliceOrMap(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Ptr:
		return m.isSliceOrMap(t.Elem())
	case reflect.Slice, reflect.Array,
		reflect.Map:
		return true
	default:
		return t == reflect.TypeOf(configure.Array{}) || t == reflect.TypeOf(configure.Map{})
	}
}

func (m *Configure) isSupportedType(t reflect.Type, depth int) bool {
	if t.Kind() == reflect.Ptr {
		return m.isSupportedType(t.Elem(), depth)
	}
	switch t.Kind() {
	case reflect.String,
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Slice, reflect.Array,
		reflect.Map:
		if depth > 1 {
			panic(configure.ExceedDeepLevel)
		}
		return m.isSupportedType(t.Elem(), depth+1)
	default:
		t = reflect.New(t).Type()
		if t.Implements(dynamicType) {
			if depth > 0 {
				panic(configure.ExceedDeepLevel)
			}
			return true
		}
		return t.Implements(staticType)
	}
}

func (m *Configure) getConfig(basePath string, cfg reflect.Value, forUpdate bool) (uint64, error) {
	pairs, err := configure.BackendImplementor().Get(basePath, true)
	if err != nil {
		return 0, err
	}
	for _, pair := range pairs.KVs {
		if err = m.setConfig(basePath, pair, cfg, forUpdate); err != nil {
			return 0, err
		}
	}
	return pairs.Version, nil
}

func (m *Configure) setConfig(basePath string, pair *configure.KVPair, cfg reflect.Value, forUpdate bool) error {
	var item ConfigItem
	if err := json.Unmarshal([]byte(pair.Value), &item); err != nil {
		return err
	}

	fieldPath := strings.Split(strings.TrimPrefix(pair.Key, basePath), "/")
	for depth := 0; depth < len(fieldPath); depth++ {
		cfg = cfg.FieldByName(fieldPath[depth])
		if !cfg.IsValid() {
			logger.Warn("substratum configure field not found", logger.Content{
				"key": pair.Key,
			})
			return nil
		}
	}
	// Try DynamicValue.
	if m.updateDynamicValue(item.Value, cfg) {
		return nil
	}
	// Try StaticValue.
	if forUpdate || m.setStaticValue(item.Value, cfg, false) {
		return nil
	}
	logger.Warn("substratum configure not updated", logger.Content{
		"key":   pair.Key,
		"value": pair.Value,
	})
	return nil
}

func (m *Configure) updateDynamicValue(s string, v reflect.Value) bool {
	if v.CanInterface() {
		if v.Type().Kind() == reflect.Ptr && v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if u, ok := v.Interface().(configure.DynamicType); ok {
			u.AtomicUpdate(s)
			return true
		}
	}
	if v.CanAddr() && v.Addr().CanInterface() {
		if u, ok := v.Addr().Interface().(configure.DynamicType); ok {
			u.AtomicUpdate(s)
			return true
		}
	}
	return false
}

func (m *Configure) setStaticValue(s string, v reflect.Value, recursion bool) bool {
	if v.CanInterface() {
		if v.Type().Kind() == reflect.Ptr && v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if u, ok := v.Interface().(configure.StaticType); ok {
			u.Set(s)
			return true
		}
	}
	if v.CanAddr() && v.Addr().CanInterface() {
		if u, ok := v.Addr().Interface().(configure.StaticType); ok {
			u.Set(s)
			return true
		}
	}
	return m.setSystemTypeValue(s, v, recursion)
}

func (m *Configure) setSystemTypeValue(s string, v reflect.Value, recursion bool) bool {
	// Used for slice or map value.
	sep := ";"
	if recursion {
		sep = ","
	}

	switch v.Type().Kind() {
	case reflect.Ptr:
		e := reflect.New(v.Type().Elem())
		m.setStaticValue(s, e.Elem(), false)
		v.Set(e)
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		bv, _ := strconv.ParseBool(s)
		v.SetBool(bv)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Type() == reflect.TypeOf(time.Second) {
			dur, _ := time.ParseDuration(s)
			v.SetInt(int64(dur))
		} else {
			iv, _ := strconv.ParseInt(s, 10, 64)
			v.SetInt(iv)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uv, _ := strconv.ParseUint(s, 10, 64)
		v.SetUint(uv)
	case reflect.Float32, reflect.Float64:
		fv, _ := strconv.ParseFloat(s, 64)
		v.SetFloat(fv)
	case reflect.Slice, reflect.Array:
		var fields []string
		if s != "" {
			fields = strings.Split(s, sep)
		}
		sv := reflect.MakeSlice(v.Type(), len(fields), len(fields))
		for i, field := range fields {
			m.setStaticValue(field, sv.Index(i), true)
		}
		v.Set(sv)
	case reflect.Map:
		var vs []string
		if s != "" {
			vs = strings.Split(s, sep)
		}
		mv := reflect.MakeMapWithSize(v.Type(), len(vs))
		for _, vv := range vs {
			kv := strings.SplitN(vv, ":", 2)
			k := reflect.New(v.Type().Key())
			nv := reflect.New(v.Type().Elem())
			m.setStaticValue(kv[0], k.Elem(), true)
			if len(kv) > 1 {
				m.setStaticValue(kv[1], nv.Elem(), true)
			}
			mv.SetMapIndex(k.Elem(), nv.Elem())
		}
		v.Set(mv)
	default:
		return false
	}
	return true
}

func (m *Configure) watchEvent(path string, ch configure.EventChan, cfg reflect.Value) {
	var (
		err error
	)

	for {
		select {
		case <-ictx.Context.Done():
			configure.BackendImplementor().Close()
			return

		case evt := <-ch:
			if evt.EventType == configure.Refresh {
				_, err = m.getConfig(path, cfg, true)
			} else {
				err = m.setConfig(path, &evt.KVPair, cfg, true)
			}
			if err != nil {
				logger.Error("substratum configure event failed", logger.Content{
					"error": err.Error(),
					"event": evt.EventType,
					"key":   evt.Key,
				})
			}
		}
	}
}

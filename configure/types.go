package configure

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	ExceedDeepLevel = "substratum: only support two level map/array"
)

// DynamicType interface.
type DynamicType interface {
	// AtomicUpdate updates value.
	AtomicUpdate(v string)

	// Changed will be invoked if value updated.
	Changed(evt UpdateEvent)
}

// StaticType interface.
type StaticType interface {
	// Set value.
	Set(v string)
}

type embedString struct {
	v atomic.Value
}

func newEmbedString(v string) *embedString {
	es := &embedString{}
	es.v.Store(v)
	return es
}

func (t *embedString) String() string {
	v := t.v.Load()
	if v == nil {
		return ""
	}
	return v.(string)
}

type String struct {
	embedString
}

func (t *String) AtomicUpdate(v string) {
	if t.String() == v {
		return
	}
	t.v.Store(v)
	callbackImpl.EvtChan() <- t
}

func (t *String) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

type embedBool struct {
	v int32
}

func (t *embedBool) String() string {
	return strconv.FormatBool(t.Bool())
}

func (t *embedBool) Bool() bool {
	return atomic.LoadInt32(&t.v) == 1
}

type Bool struct {
	embedBool
}

func (t *Bool) AtomicUpdate(v string) {
	b, _ := strconv.ParseBool(v)
	if t.Bool() == b {
		return
	}
	if b {
		atomic.StoreInt32(&t.v, 1)
	} else {
		atomic.StoreInt32(&t.v, 0)
	}
	callbackImpl.EvtChan() <- t
}

func (t *Bool) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

type embedInt struct {
	v int64
}

func (t *embedInt) String() string {
	return strconv.FormatInt(t.Int64(), 10)
}

func (t *embedInt) Int() int {
	return int(t.Int64())
}

func (t *embedInt) Int8() int8 {
	return int8(t.Int64())
}

func (t *embedInt) Int16() int16 {
	return int16(t.Int64())
}

func (t *embedInt) Int32() int32 {
	return int32(t.Int64())
}

func (t *embedInt) Int64() int64 {
	return atomic.LoadInt64(&t.v)
}

type Int struct {
	embedInt
}

func (t *Int) AtomicUpdate(v string) {
	iv, _ := strconv.ParseInt(v, 10, 64)
	if t.Int64() == iv {
		return
	}
	atomic.StoreInt64(&t.v, iv)
	callbackImpl.EvtChan() <- t
}

func (t *Int) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

type embedUint struct {
	v uint64
}

func (t *embedUint) String() string {
	return strconv.FormatUint(t.Uint64(), 10)
}

func (t *embedUint) Uint() uint {
	return uint(t.Uint64())
}

func (t *embedUint) Uint8() uint8 {
	return uint8(t.Uint64())
}

func (t *embedUint) Uint16() uint16 {
	return uint16(t.Uint64())
}

func (t *embedUint) Uint32() uint32 {
	return uint32(t.Uint64())
}

func (t *embedUint) Uint64() uint64 {
	return atomic.LoadUint64(&t.v)
}

type Uint struct {
	embedUint
}

func (t *Uint) AtomicUpdate(v string) {
	uv, _ := strconv.ParseUint(v, 10, 64)
	if t.Uint64() == uv {
		return
	}
	atomic.StoreUint64(&t.v, uv)
	callbackImpl.EvtChan() <- t
}

func (t *Uint) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

type embedFloat struct {
	v atomic.Value
}

func newEmbedFloat(f float64) *embedFloat {
	ef := &embedFloat{}
	ef.v.Store(f)
	return ef
}

func (t *embedFloat) String() string {
	return strconv.FormatFloat(t.Float64(), 'f', 6, 64)
}

func (t *embedFloat) Float32() float32 {
	return float32(t.Float64())
}

func (t *embedFloat) Float64() float64 {
	v := t.v.Load()
	if v == nil {
		return 0.0
	}
	return v.(float64)
}

type Float struct {
	embedFloat
}

func (t *Float) AtomicUpdate(v string) {
	fv, _ := strconv.ParseFloat(v, 64)
	if big.NewFloat(t.Float64()).Cmp(big.NewFloat(fv)) == 0 {
		return
	}
	t.v.Store(fv)
	callbackImpl.EvtChan() <- t
}

func (t *Float) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

type embedArray struct {
	v atomic.Value
	r bool
}

func newEmbedArray(sv []string, recursion bool) *embedArray {
	es := &embedArray{
		r: recursion,
	}
	es.v.Store(sv)
	return es
}

func (t *embedArray) load() []string {
	sv := t.v.Load()
	if sv == nil {
		return []string{}
	}
	return sv.([]string)
}

func (t *embedArray) Len() int {
	return len(t.load())
}

func (t *embedArray) String() string {
	if !t.r {
		return strings.Join(t.load(), ";")
	}
	return strings.Join(t.load(), ",")
}

func (t *embedArray) Strings() []*embedString {
	sv := t.load()
	es := make([]*embedString, 0, len(sv))
	for _, v := range sv {
		es = append(es, newEmbedString(v))
	}
	return es
}

func (t *embedArray) Bools() []*embedBool {
	sv := t.load()
	eb := make([]*embedBool, 0, len(sv))
	for _, v := range sv {
		bv := 0
		if b, _ := strconv.ParseBool(v); b {
			bv = 1
		}
		eb = append(eb, &embedBool{v: int32(bv)})
	}
	return eb
}

func (t *embedArray) Ints() []*embedInt {
	sv := t.load()
	ei := make([]*embedInt, 0, len(sv))
	for _, v := range sv {
		i, _ := strconv.ParseInt(v, 10, 64)
		ei = append(ei, &embedInt{v: i})
	}
	return ei
}

func (t *embedArray) Uints() []*embedUint {
	sv := t.load()
	eu := make([]*embedUint, 0, len(sv))
	for _, v := range sv {
		u, _ := strconv.ParseUint(v, 10, 64)
		eu = append(eu, &embedUint{v: u})
	}
	return eu
}

func (t *embedArray) Floats() []*embedFloat {
	sv := t.load()
	ef := make([]*embedFloat, 0, len(sv))
	for _, v := range sv {
		f, _ := strconv.ParseFloat(v, 64)
		ef = append(ef, newEmbedFloat(f))
	}
	return ef
}

func (t *embedArray) ArrayAt(i int) *embedArray {
	if t.r {
		panic(ExceedDeepLevel)
	}
	sv := t.load()
	if i+1 > len(sv) {
		panic("substratum: index out of range")
	}
	return newEmbedArray(strings.Split(sv[i], ","), true)
}

type Array struct {
	embedArray
}

func (t *Array) AtomicUpdate(v string) {
	if t.String() == v {
		return
	}
	sv := strings.Split(v, ";")
	t.v.Store(sv)
	callbackImpl.EvtChan() <- t
}

func (t *Array) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

type embedMap struct {
	v atomic.Value
	r bool
}

func newEmbedMap(mv map[string]string, recursion bool) *embedMap {
	em := &embedMap{
		r: recursion,
	}
	em.v.Store(mv)
	return em
}

func (t *embedMap) load() map[string]string {
	mv := t.v.Load()
	if mv == nil {
		return map[string]string{}
	} else {
		return mv.(map[string]string)
	}
}

func (t *embedMap) String() string {
	mv := t.load()
	s := make([]string, 0, len(mv))
	for k, v := range mv {
		if v == "" {
			s = append(s, k)
		} else {
			s = append(s, fmt.Sprintf("%s:%s", k, v))
		}
	}
	if !t.r {
		return strings.Join(s, ";")
	}
	return strings.Join(s, ",")
}

func (t *embedMap) Len() int {
	mv := t.load()
	return len(mv)
}

func (t *embedMap) Keys() *embedArray {
	mv := t.load()
	keys := make([]string, 0, len(mv))
	for k := range mv {
		keys = append(keys, k)
	}
	return newEmbedArray(keys, true)
}

func (t *embedMap) HasKey(key string) bool {
	mv := t.load()
	_, ok := mv[key]
	return ok
}

func (t *embedMap) StringVal(key string) *embedString {
	mv := t.load()
	return newEmbedString(mv[key])
}

func (t *embedMap) BoolVal(key string) *embedBool {
	mv := t.load()
	if v, ok := mv[key]; ok {
		if b, _ := strconv.ParseBool(v); b {
			return &embedBool{v: 1}
		}
	}
	return &embedBool{}
}

func (t *embedMap) IntVal(key string) *embedInt {
	mv := t.load()
	if v, ok := mv[key]; ok {
		i, _ := strconv.ParseInt(v, 10, 64)
		return &embedInt{v: i}
	}
	return &embedInt{}
}

func (t *embedMap) UintVal(key string) *embedUint {
	mv := t.load()
	if v, ok := mv[key]; ok {
		u, _ := strconv.ParseUint(v, 10, 64)
		return &embedUint{v: u}
	}
	return &embedUint{}
}

func (t *embedMap) FloatVal(key string) *embedFloat {
	mv := t.load()
	if v, ok := mv[key]; ok {
		f, _ := strconv.ParseFloat(v, 64)
		return newEmbedFloat(f)
	}
	return &embedFloat{}
}

func (t *embedMap) ArrayVal(key string) *embedArray {
	if t.r {
		panic(ExceedDeepLevel)
	}
	mv := t.load()
	return newEmbedArray(strings.Split(mv[key], ","), true)
}

func (t *embedMap) parse(v, sep string) map[string]string {
	vv := strings.Split(v, sep)
	mv := make(map[string]string, len(vv))
	for _, v := range vv {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) == 1 {
			mv[parts[0]] = ""
		} else {
			mv[parts[0]] = parts[1]
		}
	}
	return mv
}

func (t *embedMap) MapVal(key string) *embedMap {
	if t.r {
		panic(ExceedDeepLevel)
	}
	mv := t.load()
	m := t.parse(mv[key], ",")
	return newEmbedMap(m, true)
}

type Map struct {
	embedMap
}

func (t *Map) AtomicUpdate(v string) {
	if t.String() == v {
		return
	}
	mv := t.parse(v, ";")
	t.v.Store(mv)
	callbackImpl.EvtChan() <- t
}

func (t *Map) Changed(evt UpdateEvent) {
	callbackImpl.RegChan() <- &CallbackFunc{
		Value: t,
		Event: evt,
	}
}

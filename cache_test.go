package sessions

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestGetCacheFromRequestContext(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	getCacheFromRequestContext(r)
}

func TestPut(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	s.Put(r, "foo", "bar")

	if c.Data["foo"] != "bar" {
		t.Errorf("got %q: expected %q", c.Data["foo"], "bar")
	}

	if !c.modified {
		t.Errorf("got %v: expected %v", c.modified, true)
	}
}

func TestGet(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	str, ok := s.Get(r, "foo").(string)
	if !ok {
		t.Errorf("could not convert %T to string", s.Get(r, "foo"))
	}

	if str != "bar" {
		t.Errorf("got %q: expected %q", str, "bar")
	}
}

func TestPop(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	str, ok := s.Pop(r, "foo").(string)
	if !ok {
		t.Errorf("could not convert %T to string", s.Get(r, "foo"))
	}

	if str != "bar" {
		t.Errorf("got %q: expected %q", str, "bar")
	}

	_, ok = c.Data["foo"]
	if ok {
		t.Errorf("got %v: expected %v", ok, false)
	}

	if !c.modified {
		t.Errorf("got %v: expected %v", c.modified, true)
	}
}

func TestRemove(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	s.Remove(r, "foo")

	if c.Data["foo"] != nil {
		t.Errorf("got %v: expected %v", c.Data["foo"], nil)
	}

	if !c.modified {
		t.Errorf("got %v: expected %v", c.modified, true)
	}
}

func TestExists(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	if !s.Exists(r, "foo") {
		t.Errorf("got %v: expected %v", s.Exists(r, "foo"), true)
	}

	if s.Exists(r, "baz") {
		t.Errorf("got %v: expected %v", s.Exists(r, "baz"), false)
	}
}

func TestKeys(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	c.Data["woo"] = "waa"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	keys := s.Keys(r)
	if !reflect.DeepEqual(keys, []string{"foo", "woo"}) {
		t.Errorf("got %v: expected %v", keys, []string{"foo", "woo"})
	}
}

func TestGetString(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	str := s.GetString(r, "foo")
	if str != "bar" {
		t.Errorf("got %q: expected %q", str, "bar")
	}

	str = s.GetString(r, "baz")
	if str != "" {
		t.Errorf("got %q: expected %q", str, "")
	}
}

func TestGetBool(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = true
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	b := s.GetBool(r, "foo")
	if b != true {
		t.Errorf("got %v: expected %v", b, true)
	}

	b = s.GetBool(r, "baz")
	if b != false {
		t.Errorf("got %v: expected %v", b, false)
	}
}

func TestGetInt(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = 123
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	i := s.GetInt(r, "foo")
	if i != 123 {
		t.Errorf("got %v: expected %d", i, 123)
	}

	i = s.GetInt(r, "baz")
	if i != 0 {
		t.Errorf("got %v: expected %d", i, 0)
	}
}

func TestGetFloat(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = 123.456
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	f := s.GetFloat(r, "foo")
	if f != 123.456 {
		t.Errorf("got %v: expected %f", f, 123.456)
	}

	f = s.GetFloat(r, "baz")
	if f != 0 {
		t.Errorf("got %v: expected %f", f, 0.00)
	}
}

func TestGetBytes(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = []byte("bar")
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	b := s.GetBytes(r, "foo")
	if !bytes.Equal(b, []byte("bar")) {
		t.Errorf("got %v: expected %v", b, []byte("bar"))
	}

	b = s.GetBytes(r, "baz")
	if b != nil {
		t.Errorf("got %v: expected %v", b, nil)
	}
}

func TestGetTime(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()

	c := newCache(time.Hour)
	c.Data["foo"] = now
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	tm := s.GetTime(r, "foo")
	if tm != now {
		t.Errorf("got %v: expected %v", tm, now)
	}

	tm = s.GetTime(r, "baz")
	if !tm.IsZero() {
		t.Errorf("got %v: expected %v", tm, time.Time{})
	}
}

func TestPopString(t *testing.T) {
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := newCache(time.Hour)
	c.Data["foo"] = "bar"
	r = addCacheToRequestContext(r, c)

	s := New([]byte("secret"))
	str := s.PopString(r, "foo")
	if str != "bar" {
		t.Errorf("got %q: expected %q", str, "bar")
	}

	_, ok := c.Data["foo"]
	if ok {
		t.Errorf("got %v: expected %v", ok, false)
	}

	if !c.modified {
		t.Errorf("got %v: expected %v", c.modified, true)
	}

	str = s.PopString(r, "bar")
	if str != "" {
		t.Errorf("got %q: expected %q", str, "")
	}
}

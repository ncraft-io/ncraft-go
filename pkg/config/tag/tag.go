package tag

import (
    "flag"
    "fmt"
    "reflect"
    "strconv"
    "strings"
    "time"

    "github.com/fatih/structs"
)

// Loader satisfies the loader interface. It parses a struct's field tags
// and populates the each field with that given tag.
type Loader struct {
    // DefaultTagName is the default tag name for struct fields to define
    // default values for a field. Example:
    //
    //   // Field's default value is "koding".
    //   Name string `default:"koding"`
    //
    // The default value is "default" if it's not set explicitly.
    DefaultTagName string
}

func (t *Loader) Load(s interface{}) error {
    if t.DefaultTagName == "" {
        t.DefaultTagName = "default"
    }

    for _, field := range structs.Fields(s) {

        if err := t.processField(t.DefaultTagName, field); err != nil {
            return err
        }
    }

    return nil
}

// processField gets tagName and the field, recursively checks if the field has the given
// tag, if yes, sets it otherwise ignores
func (t *Loader) processField(tagName string, field *structs.Field) error {
    switch field.Kind() {
    case reflect.Struct:
        for _, f := range field.Fields() {
            if err := t.processField(tagName, f); err != nil {
                return err
            }
        }
    default:
        defaultVal := field.Tag(t.DefaultTagName)
        if defaultVal == "" {
            return nil
        }

        err := fieldSet(field, defaultVal)
        if err != nil {
            return err
        }
    }

    return nil
}

// fieldSet sets field value from the given string value. It converts the
// string value in a sane way and is usefulf or environment variables or flags
// which are by nature in string types.
func fieldSet(field *structs.Field, v string) error {
    switch f := field.Value().(type) {
    case flag.Value:
        if v := reflect.ValueOf(field.Value()); v.IsNil() {
            typ := v.Type()
            if typ.Kind() == reflect.Ptr {
                typ = typ.Elem()
            }

            if err := field.Set(reflect.New(typ).Interface()); err != nil {
                return err
            }

            f = field.Value().(flag.Value)
        }

        return f.Set(v)
    }

    // TODO: add support for other types
    switch field.Kind() {
    case reflect.Bool:
        val, err := strconv.ParseBool(v)
        if err != nil {
            return err
        }

        if err := field.Set(val); err != nil {
            return err
        }
    case reflect.Int:
        i, err := strconv.Atoi(v)
        if err != nil {
            return err
        }

        if err := field.Set(i); err != nil {
            return err
        }
    case reflect.String:
        if err := field.Set(v); err != nil {
            return err
        }
    case reflect.Slice:
        switch t := field.Value().(type) {
        case []string:
            if err := field.Set(strings.Split(v, ",")); err != nil {
                return err
            }
        case []int:
            var list []int
            for _, in := range strings.Split(v, ",") {
                i, err := strconv.Atoi(in)
                if err != nil {
                    return err
                }

                list = append(list, i)
            }

            if err := field.Set(list); err != nil {
                return err
            }
        default:
            return fmt.Errorf("multiconfig: field '%s' of type slice is unsupported: %s (%T)",
                field.Name(), field.Kind(), t)
        }
    case reflect.Float64:
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return err
        }

        if err := field.Set(f); err != nil {
            return err
        }
    case reflect.Int64:
        switch t := field.Value().(type) {
        case time.Duration:
            d, err := time.ParseDuration(v)
            if err != nil {
                return err
            }

            if err := field.Set(d); err != nil {
                return err
            }
        case int64:
            p, err := strconv.ParseInt(v, 10, 0)
            if err != nil {
                return err
            }

            if err := field.Set(p); err != nil {
                return err
            }
        default:
            return fmt.Errorf("multiconfig: field '%s' of type int64 is unsupported: %s (%T)",
                field.Name(), field.Kind(), t)
        }

    default:
        return fmt.Errorf("multiconfig: field '%s' has unsupported type: %s", field.Name(), field.Kind())
    }

    return nil
}

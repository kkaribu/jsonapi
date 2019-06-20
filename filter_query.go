package jsonapi

import (
	"encoding/json"
	"time"
)

// A Condition is used to define filters when querying collections.
type Condition struct {
	Field string      `json:"f"`
	Op    string      `json:"o"`
	Val   interface{} `json:"v"`
	Col   string      `json:"c"`
}

// cnd is an internal version of Condition.
type cnd struct {
	Field string          `json:"f"`
	Op    string          `json:"o"`
	Val   json.RawMessage `json:"v"`
	Col   string          `json:"c"`
}

// UnmarshalJSON parses the provided data and populates a Condition.
func (c *Condition) UnmarshalJSON(data []byte) error {
	tmpCnd := cnd{}
	err := json.Unmarshal(data, &tmpCnd)
	if err != nil {
		return err
	}

	c.Field = tmpCnd.Field
	c.Op = tmpCnd.Op
	c.Col = tmpCnd.Col

	if tmpCnd.Op == "and" || tmpCnd.Op == "or" {
		c.Field = ""

		cnds := []Condition{}
		err := json.Unmarshal(tmpCnd.Val, &cnds)
		if err != nil {
			return err
		}
		c.Val = cnds
	} else if tmpCnd.Op == "=" ||
		tmpCnd.Op == "!=" ||
		tmpCnd.Op == "<" ||
		tmpCnd.Op == "<=" ||
		tmpCnd.Op == ">" ||
		tmpCnd.Op == ">=" {

		err := json.Unmarshal(tmpCnd.Val, &(c.Val)) // TODO parenthesis needed?
		if err != nil {
			return err
		}
	}

	return nil
}

// MarshalJSON marshals a Condition into JSON.
func (c *Condition) MarshalJSON() ([]byte, error) {
	payload := map[string]interface{}{}
	if c.Field != "" {
		payload["f"] = c.Field
	}
	if c.Op != "" {
		payload["o"] = c.Op
	}
	payload["v"] = c.Val
	if c.Col != "" {
		payload["c"] = c.Col
	}
	return json.Marshal(payload)
}

// FilterResource reports whether res is valid under the rules defined
// in cond.
func FilterResource(res Resource, cond *Condition) bool {
	var (
		val interface{}
		// typ string
	)
	if _, ok := res.Attrs()[cond.Field]; ok {
		val = res.Get(cond.Field)
	}
	if rel, ok := res.Rels()[cond.Field]; ok {
		if rel.ToOne {
			val = res.GetToOne(cond.Field)
		} else {
			val = res.GetToMany(cond.Field)
		}
	}

	switch cond.Op {
	case "and":
		conds := cond.Val.([]*Condition)
		for i := range conds {
			if !FilterResource(res, conds[i]) {
				return false
			}
		}
	case "or":
		conds := cond.Val.([]*Condition)
		for i := range conds {
			if FilterResource(res, conds[i]) {
				return true
			}
		}
	case "=", "!=", "<", "<=", ">", ">=":
		if val == cond.Val {
			return checkVal(cond.Op, val, cond.Val)
		}
	}

	return false
}

func checkVal(op string, rval, cval interface{}) bool {
	switch rval.(type) {
	case string:
		return checkStr(op, rval.(string), cval.(string))
	case int:
		return checkInt(op, int64(rval.(int)), int64(cval.(int)))
	case int8:
		return checkInt(op, int64(rval.(int8)), int64(cval.(int8)))
	case int16:
		return checkInt(op, int64(rval.(int16)), int64(cval.(int16)))
	case int32:
		return checkInt(op, int64(rval.(int32)), int64(cval.(int32)))
	case int64:
		return checkInt(op, rval.(int64), cval.(int64))
	case uint:
		return checkUint(op, uint64(rval.(uint)), uint64(cval.(uint)))
	case uint8:
		return checkUint(op, uint64(rval.(uint8)), uint64(cval.(uint8)))
	case uint16:
		return checkUint(op, uint64(rval.(uint16)), uint64(cval.(uint16)))
	case uint32:
		return checkUint(op, uint64(rval.(uint32)), uint64(cval.(uint32)))
	case uint64:
		return checkUint(op, rval.(uint64), cval.(uint64))
	case bool:
		return checkBool(op, rval.(bool), cval.(bool))
	case time.Time:
		return checkTime(op, rval.(time.Time), cval.(time.Time))
	default:
		return false
	}
}

func checkStr(op string, rval, cval string) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	case "<":
		return rval < cval
	case "<=":
		return rval <= cval
	case ">":
		return rval > cval
	case ">=":
		return rval >= cval
	default:
		return false
	}
}

func checkInt(op string, rval, cval int64) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	case "<":
		return rval < cval
	case "<=":
		return rval <= cval
	case ">":
		return rval > cval
	case ">=":
		return rval >= cval
	default:
		return false
	}
}

func checkUint(op string, rval, cval uint64) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	case "<":
		return rval < cval
	case "<=":
		return rval <= cval
	case ">":
		return rval > cval
	case ">=":
		return rval >= cval
	default:
		return false
	}
}

func checkBool(op string, rval, cval bool) bool {
	switch op {
	case "=":
		return rval == cval
	case "!=":
		return rval != cval
	default:
		return false
	}
}

func checkTime(op string, rval, cval time.Time) bool {
	switch op {
	case "=":
		return rval.Equal(cval)
	case "!=":
		return !rval.Equal(cval)
	case "<":
		return rval.Before(cval)
	case "<=":
		return rval.Before(cval) || rval.Equal(cval)
	case ">":
		return rval.After(cval) || rval.Equal(cval)
	case ">=":
		return rval.After(cval) || rval.Equal(cval)
	default:
		return false
	}
}

func makeInt(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	default:
		return 0, false
	}
}

func makeUint(val interface{}) (uint64, bool) {
	switch v := val.(type) {
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	default:
		return 0, false
	}
}

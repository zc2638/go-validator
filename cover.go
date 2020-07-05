/**
 * Created by zc on 2020/7/5.
 */
package validator

import "encoding/json"

func JSONCover() Cover {
	return func(data []byte, s interface{}) error {
		return json.Unmarshal(data, s)
	}
}
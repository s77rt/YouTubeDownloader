package youtube

import (
	"fmt"
	"encoding/json"
	"github.com/kkdai/youtube/v2"
)

func assertVideo(video interface{}) (*youtube.Video, error) {
	v, ok := video.(*youtube.Video)
	if ok == true {
		return v, nil
	} else {
		switch vv := video.(type) {
		case map[string]interface {}:
			v_json, err := json.Marshal(vv)
			if err != nil {
				return nil, err
			}

			var new_v *youtube.Video
			err = json.Unmarshal(v_json, &new_v)
			if err != nil {
				return nil, err
			}

			return new_v, nil
		default:
			return nil, fmt.Errorf("not a video object")
		}
	}
}

func assertFormat(format interface{}) (*youtube.Format, error) {
	f, ok := format.(*youtube.Format)
	if ok == true {
		return f, nil
	} else {
		switch ff := format.(type) {
		case map[string]interface {}:
			f_json, err := json.Marshal(ff)
			if err != nil {
				return nil, err
			}

			var new_f *youtube.Format
			err = json.Unmarshal(f_json, &new_f)
			if err != nil {
				return nil, err
			}

			return new_f, nil
		default:
			return nil, fmt.Errorf("not a format object")
		}
	}
}

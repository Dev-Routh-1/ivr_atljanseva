package asterisk

import "fmt"

type Vars map[string]string

const (
	Phone           = "IVR_PHONE"
	Language        = "IVR_LANGUAGE"
	Pincode         = "IVR_PINCODE"
	Ward            = "IVR_WARD"
	NagarsevakID    = "IVR_NAGARSEVAK_ID"
	NagarsevakName  = "IVR_NAGARSEVAK_NAME"
	Found           = "IVR_FOUND"
	Status          = "IVR_STATUS"
)

func (v Vars) AGI() []string {
	out := make([]string, 0, len(v))
	for name, val := range v {
		out = append(out, fmt.Sprintf("SET VARIABLE %s %q", name, val))
	}
	return out
}

func FromRequest(phone, language, pincode, ward string) Vars {
	return Vars{
		Phone:    phone,
		Language: language,
		Pincode:  pincode,
		Ward:     ward,
	}
}

func FromCitizenLookup(found bool, language, pincode, ward, nagarsevakID, nagarsevakName string) Vars {
	v := Vars{
		Found:    fmt.Sprintf("%t", found),
		Language: language,
		Pincode:  pincode,
		Ward:     ward,
	}
	if nagarsevakID != "" {
		v[NagarsevakID] = nagarsevakID
	}
	if nagarsevakName != "" {
		v[NagarsevakName] = nagarsevakName
	}
	return v
}

func FromResolveWard(status, phone, language, pincode, ward, nagarsevakID, nagarsevakName string) Vars {
	v := Vars{
		Status:   status,
		Phone:    phone,
		Language: language,
		Pincode:  pincode,
	}
	if ward != "" {
		v[Ward] = ward
	}
	if nagarsevakID != "" {
		v[NagarsevakID] = nagarsevakID
	}
	if nagarsevakName != "" {
		v[NagarsevakName] = nagarsevakName
	}
	return v
}

func FromNagarsevak(status, phone, language, pincode, ward, nagarsevakID, nagarsevakName string) Vars {
	v := Vars{
		Status:   status,
		Phone:    phone,
		Language: language,
		Pincode:  pincode,
		Ward:     ward,
	}
	if nagarsevakID != "" {
		v[NagarsevakID] = nagarsevakID
	}
	if nagarsevakName != "" {
		v[NagarsevakName] = nagarsevakName
	}
	return v
}

func FromComplete(phone, language, pincode, ward, nagarsevakName string) Vars {
	return Vars{
		Status:         "saved",
		Phone:          phone,
		Language:       language,
		Pincode:        pincode,
		Ward:           ward,
		NagarsevakName: nagarsevakName,
	}
}

package env

import (
	"fmt"
	"strings"
)

// envSet represents a set of environment variables.
type envSet map[string]string

// Override represents a set of environment variables changes, corresponding to
// os.Setenv and os.Unsetenv operations.
type Override map[string]*string

// Apply applies a Override to EnvSet, modifying its contents.
func (e envSet) Apply(orr Override, v interface{}) (envSet, error) {
	for k, v := range orr {
		if v == nil {
			// Equivalent to os.Unsetenv
			delete(e, k)
		} else {
			// Equivalent to os.Setenv
			e[k] = *v
		}
	}
	return e, unmarshal(e, v)
}

// envToEnvSet transforms a slice of string with the format "key=value" into
// the corresponding EnvSet. If any item in environ does follow the format,
// envToEnvSet returns ErrInvalidEnviron.
func envToEnvSet(environ []string) (envSet, error) {
	m := make(envSet)
	for _, v := range environ {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			return nil, errInvalidEnviron
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}

// envSetToEnv transforms a EnvSet into a slice of strings with the format
// "key=value".
func envSetToEnv(m envSet) []string {
	var environ []string
	for k, v := range m {
		environ = append(environ, fmt.Sprintf("%s=%s", k, v))
	}
	return environ
}

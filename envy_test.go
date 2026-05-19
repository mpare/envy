package envy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// String constants
	defaultAppName    = "app"
	defaultHost       = "localhost"
	defaultDBHost     = "db.example.com"
	defaultDBPassword = "secret123"
	invalidNumber     = "not-a-number"
	invalidBool       = "maybe"
	invalidDuration   = "not-a-duration"
	invalidFloat      = "not-a-float"
	invalidPointer    = "not-a-pointer"

	// Boolean string values
	boolTrue  = "true"
	boolFalse = "false"
	boolOne   = "1"
	boolZero  = "0"
	boolT     = "t"
	boolF     = "f"

	// Duration/Port string values
	timeout30s    = "30s"
	customPortStr = "9090"

	// Numeric constants
	customPort           = 9090
	defaultDBPort        = 5432
	expectedIPCount      = 3
	expectedPortCount    = 3
	port8081             = 8081
	port8082             = 8082
	firstPort            = 8080
	customTimeoutSeconds = 30

	// IP addresses
	ip1 = "192.168.1.1"
	ip2 = "192.168.1.2"
	ip3 = "192.168.1.3"

	// Tag values
	tag1    = "web"
	tag2    = "api"
	tag3    = "v2"
	tagCSV  = "web,api,v2"
	ipCSV   = "192.168.1.1;192.168.1.2;192.168.1.3"
	portCSV = "8080,8081,8082"
)

func Test_givenAllTypes_whenLoadFrom_thenAllFieldsPopulated(t *testing.T) {
	// given
	type Config struct {
		Name    string        `env:"NAME,default=app"`
		Port    int           `env:"PORT,default=8080"`
		Debug   bool          `env:"DEBUG,default=false"`
		Timeout time.Duration `env:"TIMEOUT,default=5s"`
		Rate    float64       `env:"RATE,default=1.0"`
		Tags    []string      `env:"TAGS,separator=,"`
	}

	var cfg Config
	envMap := map[string]string{
		"PORT":    customPortStr,
		"DEBUG":   boolTrue,
		"TIMEOUT": timeout30s,
		"TAGS":    tagCSV,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.NoError(t, err)
	assert.Equal(t, defaultAppName, cfg.Name)
	assert.Equal(t, customPort, cfg.Port)
	assert.True(t, cfg.Debug)
	assert.Equal(t, customTimeoutSeconds*time.Second, cfg.Timeout)
	assert.Equal(t, 1.0, cfg.Rate)
	assert.ElementsMatch(t, []string{tag1, tag2, tag3}, cfg.Tags)
}

func Test_givenRequiredField_whenMissing_thenValidationError(t *testing.T) {
	// given
	type Config struct {
		Secret string `env:"SECRET,required"`
	}

	var cfg Config
	envMap := map[string]string{}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.Error(t, err)
	validationError, ok := err.(*ValidationError)
	require.True(t, ok)
	assert.Equal(t, 1, len(validationError.Errors))
}

func Test_givenDefaultValues_whenLoadFrom_thenDefaultsApplied(t *testing.T) {
	// given
	type Config struct {
		Host string `env:"HOST,default=localhost"`
		Port int    `env:"PORT,default=5432"`
	}

	var cfg Config
	envMap := map[string]string{}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.NoError(t, err)
	assert.Equal(t, defaultHost, cfg.Host)
	assert.Equal(t, defaultDBPort, cfg.Port)
}

func Test_givenNestedStruct_whenLoadFromWithPrefix_thenNestedFieldsPopulated(t *testing.T) {
	// given
	type DatabaseConfig struct {
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT,default=5432"`
		Password string `env:"PASSWORD,required"`
	}

	type Config struct {
		AppName string         `env:"APP_NAME,default=myapp"`
		DB      DatabaseConfig `env:",prefix=DB_"`
	}

	var cfg Config
	envMap := map[string]string{
		"DB_HOST":     defaultDBHost,
		"DB_PASSWORD": defaultDBPassword,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.NoError(t, err)
	assert.Equal(t, "myapp", cfg.AppName)
	assert.Equal(t, defaultDBHost, cfg.DB.Host)
	assert.Equal(t, defaultDBPort, cfg.DB.Port)
	assert.Equal(t, defaultDBPassword, cfg.DB.Password)
}

func Test_givenNestedStructRequiredField_whenMissing_thenValidationError(t *testing.T) {
	// given
	type DatabaseConfig struct {
		Host string `env:"HOST,required"`
	}

	type Config struct {
		DB DatabaseConfig `env:",prefix=DB_"`
	}

	var cfg Config
	envMap := map[string]string{}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.Error(t, err)
	_, ok := err.(*ValidationError)
	assert.True(t, ok)
}

func Test_givenInvalidInt_whenLoadFrom_thenValidationError(t *testing.T) {
	// given
	type Config struct {
		Port int `env:"PORT"`
	}

	var cfg Config
	envMap := map[string]string{
		"PORT": invalidNumber,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.Error(t, err)
	_, ok := err.(*ValidationError)
	assert.True(t, ok)
}

func Test_givenInvalidBool_whenLoadFrom_thenValidationError(t *testing.T) {
	// given
	type Config struct {
		Debug bool `env:"DEBUG"`
	}

	var cfg Config
	envMap := map[string]string{
		"DEBUG": invalidBool,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.Error(t, err)
	_, ok := err.(*ValidationError)
	assert.True(t, ok)
}

func Test_givenInvalidDuration_whenLoadFrom_thenValidationError(t *testing.T) {
	// given
	type Config struct {
		Timeout time.Duration `env:"TIMEOUT"`
	}

	var cfg Config
	envMap := map[string]string{
		"TIMEOUT": invalidDuration,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.Error(t, err)
	_, ok := err.(*ValidationError)
	assert.True(t, ok)
}

func Test_givenInvalidFloat_whenLoadFrom_thenValidationError(t *testing.T) {
	// given
	type Config struct {
		Rate float64 `env:"RATE"`
	}

	var cfg Config
	envMap := map[string]string{
		"RATE": invalidFloat,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.Error(t, err)
	_, ok := err.(*ValidationError)
	assert.True(t, ok)
}

func Test_givenSliceWithCustomSeparator_whenLoadFrom_thenSlicePopulated(t *testing.T) {
	// given
	type Config struct {
		IPs []string `env:"IPS,separator=;"`
	}

	var cfg Config
	envMap := map[string]string{
		"IPS": ipCSV,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedIPCount, len(cfg.IPs))
	assert.ElementsMatch(t, []string{ip1, ip2, ip3}, cfg.IPs)
}

func Test_givenIntSlice_whenLoadFrom_thenSlicePopulated(t *testing.T) {
	// given
	type Config struct {
		Ports []int `env:"PORTS,separator=,"`
	}

	var cfg Config
	envMap := map[string]string{
		"PORTS": portCSV,
	}

	// when
	err := LoadFrom(&cfg, envMap)

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedPortCount, len(cfg.Ports))
	assert.ElementsMatch(t, []int{firstPort, port8081, port8082}, cfg.Ports)
}

func Test_givenBoolVariants_whenLoadFrom_thenCorrectValueDecoded(t *testing.T) {
	testCases := []struct {
		name     string
		value    string
		expected bool
	}{
		{"true_string", boolTrue, true},
		{"false_string", boolFalse, false},
		{"one", boolOne, true},
		{"zero", boolZero, false},
		{"t_lowercase", boolT, true},
		{"f_lowercase", boolF, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			type Config struct {
				Value bool `env:"VALUE"`
			}

			var cfg Config
			envMap := map[string]string{
				"VALUE": testCase.value,
			}

			// when
			err := LoadFrom(&cfg, envMap)

			// then
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, cfg.Value)
		})
	}
}

func Test_givenInvalidPointer_whenLoadFrom_thenError(t *testing.T) {
	// given / when
	err := LoadFrom(invalidPointer, map[string]string{})

	// then
	assert.Error(t, err)
}

func Test_givenNilPointer_whenLoadFrom_thenError(t *testing.T) {
	// given
	var cfg *struct{}
	cfg = nil

	// when
	err := LoadFrom(cfg, map[string]string{})

	// then
	assert.Error(t, err)
}

func Test_givenPointerToNonStruct_whenLoadFrom_thenError(t *testing.T) {
	// given
	var str string

	// when
	err := LoadFrom(&str, map[string]string{})

	// then
	assert.Error(t, err)
}

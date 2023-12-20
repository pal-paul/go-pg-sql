# env

Package env provides an env struct field tag to marshal and unmarshal environment variables

## Usage

```go

type Environment struct {
	ProjectId   string `env:"PROJECT_ID"`
	LandingPage struct {
		Url       string `env:"LANDING_PAGE_BASE_URL"`
		OptOutUrl string `env:"OPTOUT_PAGE_BASE_URL"`
	}
	DefaultValue string `env:"MISSING_VAR,default=default_value"`
}

func main() {
	var environment Environment
	e, err := env.Unmarshal(&environment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(environment.LandingPage.Url)

	url := "new_value"
	cs := env.Override{
		"LANDING_PAGE_BASE_URL": &url,
	}
	_, err = e.Apply(cs, &environment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(environment.LandingPage.Url)
}
```
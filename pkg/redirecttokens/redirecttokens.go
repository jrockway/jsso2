package redirecttokens

import (
	"fmt"
	"time"

	"github.com/jrockway/jsso2/pkg/tokens"
	"github.com/jrockway/jsso2/pkg/types"
)

const RedirectTokenLifetime = 5 * time.Minute

type Config struct {
	tokens.GeneratorConfig
}

func (c *Config) New(dest string) (string, error) {
	msg := &types.RedirectToken{
		Uri: dest,
	}
	token, err := tokens.New(msg, c.Key)
	if err != nil {
		return "", fmt.Errorf("generate redirect token: %w", err)
	}
	return token, nil
}

func (c *Config) Unmarshal(token string) (string, error) {
	msg := &types.RedirectToken{}
	if err := tokens.VerifyAndUnmarshal(msg, token, RedirectTokenLifetime, c.Key); err != nil {
		return "", fmt.Errorf("verify and unmarshal redirect token: %w", err)
	}
	return msg.GetUri(), nil
}

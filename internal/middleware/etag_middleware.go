package middleware

import (
	"fmt"

	"github.com/dchest/siphash"
	"github.com/gofiber/fiber/v2"
)

type SipHashETagMiddleware struct {
	key0 uint64
	key1 uint64
}

func NewSipHashETagMiddleware(k0, k1 uint64) *SipHashETagMiddleware {
	return &SipHashETagMiddleware{
		key0: k0,
		key1: k1,
	}
}

func (m *SipHashETagMiddleware) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if err := c.Next(); err != nil {
			return err
		}

		body := c.Response().Body()
		if len(body) == 0 {
			return nil
		}

		// strong instead of semantic equivalence
		etag := fmt.Sprintf(`"%x"`, siphash.Hash(m.key0, m.key1, body))

		if inm := c.Get("If-None-Match"); inm != "" {
			if inm == etag {
				c.Status(fiber.StatusNotModified)
				c.Response().SetBody(nil)
				return nil
			}
		}

		c.Set("ETag", etag)

		sig := siphash.Hash(m.key0, m.key1, body)
		c.Set("X-Content-Signature", fmt.Sprintf("%x", sig))

		return nil
	}
}

package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func parsePath[T any](ctx *fiber.Ctx) (*T, error) {
	return parseParamsWithStatus[T](ctx.ParamsParser)
}

func parseQuery[T any](ctx *fiber.Ctx) (*T, error) {
	return parseParamsWithStatus[T](ctx.QueryParser)
}

func parseHeaders[T any](ctx *fiber.Ctx) (*T, error) {
	return parseParamsWithStatus[T](ctx.ReqHeaderParser)
}

func parseBody[T any](ctx *fiber.Ctx) (*T, error) {
	return parseParamsWithStatus[T](ctx.BodyParser)
}

type parseParamsFunc func(out any) error

func parseParamsWithStatus[T any](
	parse parseParamsFunc,
) (*T, error) {
	params, err := parseParams[T](parse)
	if err != nil {
		return nil, fiber.NewError(
			fiber.StatusBadRequest,
			fmt.Sprintf("parse param: %v", err),
		)
	}

	return params, nil
}

func parseParams[T any](
	parseRequest parseParamsFunc,
) (*T, error) {
	var params T
	if err := parseRequest(&params); err != nil {
		return nil, err
	}

	return &params, nil
}

package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"PieceIndex",
		"GET",
		"/",
		PieceIndex,
	},
	Route{
		"PieceView",
		"GET",
		"/piece/{pieceId}",
		PieceView,
	},
        Route{
                "PieceAdd",
                "POST",
                "/piece/add",
                PieceAdd,
        },
	Route{
		"PieceSearch",
		"GET",
		"/piece/search/byname/{pieceName}",
		PieceSearch,
	},
	Route{
		"OutView",
		"GET",
		"/out/{outId}",
		OutView,
	},
	Route{
		"OutHandIn",
		"POST",
		"/out/checkin",
		OutHandIn,
	},
	Route{
		"OutHandOut",
		"POST",
		"/out/checkout",
		OutHandOut,
	},
        Route{
                "LoanRequest",
                "POST",
                "/loan/request",
                LoanSubmit,
        },
        Route{
                "LoanApprove",
                "POST",
                "/loan/approve",
                LoanApprove,
        },
}

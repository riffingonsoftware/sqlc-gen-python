package python

import (
	"log"
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func sqliteType(req *plugin.GenerateRequest, col *plugin.Column) string {
	columnType := strings.ToLower(sdk.DataType(col.Type))

	switch columnType {
	case "int", "integer", "tinyint", "smallint", "mediumint", "bigint", "unsigned big int", "int2", "int8":
		return "int"
	case "real", "double", "double precision", "float":
		return "float"
	case "numeric", "decimal":
		return "decimal.Decimal"
	case "boolean":
		return "bool"
	case "json":
		return "Any"
	case "blob":
		return "memoryview"
	case "date":
		return "datetime.date"
	case "datetime":
		return "datetime.datetime"
	case "text", "character", "varchar", "nchar", "nvarchar", "clob":
		return "str"
	default:
		for _, schema := range req.Catalog.Schemas {
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return "models." + modelName(enum.Name, req.Settings)
					}
					return "models." + modelName(schema.Name+"_"+enum.Name, req.Settings)
				}
			}
		}
		log.Printf("unknown SQLite type: %s\n", columnType)
		return "Any"
	}
}

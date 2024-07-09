package python

import (
	"log"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func sqliteType(req *plugin.GenerateRequest, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)

	switch columnType {
	case "INTEGER", "INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT", "INT2", "INT8":
		return "int"
	case "REAL", "DOUBLE", "DOUBLE PRECISION", "FLOAT":
		return "float"
	case "NUMERIC", "DECIMAL":
		return "decimal.Decimal"
	case "BOOLEAN":
		return "bool"
	case "JSON":
		return "Any"
	case "BLOB":
		return "memoryview"
	case "DATE":
		return "datetime.date"
	case "TIME":
		return "datetime.time"
	case "DATETIME", "TIMESTAMP":
		return "datetime.datetime"
	case "TEXT", "CHARACTER", "VARCHAR", "NCHAR", "NVARCHAR", "CLOB":
		return "str"
	default:
		// SQLite doesn't have built-in UUID, INET, CIDR, MACADDR, or INTERVAL types
		// It also doesn't have an ENUM type, but we'll keep the enum check for consistency
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

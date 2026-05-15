package postgres

import (
	"database/sql"
	"log_parser3000/internal/domain"
)

func nullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func nodeInfoGUID(ni *domain.NodeInfo) string {
	if ni.GeneralInfo != nil {
		return ni.GeneralInfo.NodeGUID
	}
	if ni.SwitchInfo != nil {
		return ni.SwitchInfo.NodeGUID
	}
	if ni.SharpInfo != nil {
		return ni.SharpInfo.NodeGUID
	}
	return ""
}

func nullableIntPtr(v *int) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*v), Valid: true}
}

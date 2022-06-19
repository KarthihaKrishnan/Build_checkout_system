package controllers

//	"database/sql"
//	"encoding/json"
//	"net/http"

//	"Build_checkout_system/pkg/models"
//	"Build_checkout_system/pkg/utils"

//	"github.com/jmoiron/sqlx"

// FuncMacBookProPromotion function to check availability of MacBook Pro
/*func FuncMacBookProPromotion(db *sqlx.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Declarations
		respType := utils.ResponseType{
			W: w,
			R: r,
		}
		var err error

		type Result struct {
			Status     int    `json:"status"`
			StatusText string `json:"status_text"`
		}

		type Responses struct {
			Result
			Results interface{}
		}

		type ProductResult struct {
			Scanned_Items string `json:"scanned_item"`
			Total float64 `json:"total"`
		}

		w.Header().Set("Content-Type", "application/json")
		//start processing

		product := models.Product{}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err = decoder.Decode(&product)
		if err != nil {
			utils.ErrorResponseHandler("Invalid request object supplied", http.StatusBadRequest, err, respType)
			return
		}

		var tx *sql.Tx
		tx, err = db.Begin()
		if err != nil {
			utils.ErrorResponseHandler("Internal server error 1", http.StatusInternalServerError, err, respType)
			return
		}

		if

	})
} */

package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/pcdummy/albiondata-api/lib"
	adslib "github.com/pcdummy/albiondata-sql/lib"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	cfgFile string
	db      *gorm.DB
)

var rootCmd = &cobra.Command{
	Use:   "albiondata-api",
	Short: "albiondata-api is the API Server for the Albion Data Project",
	Long: `Reads data from a SQL Database (MSSQL, MySQL, PostgreSQL and SQLite3 are supported), 
and serves them through a HTTP API`,
	Run: doCmd,
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.albiondata-api.yaml")
	rootCmd.PersistentFlags().StringP("listen", "l", "[::1]:3080", "Host and port to listen on")
	rootCmd.PersistentFlags().StringP("dbType", "t", "mysql", "Database type must be one of mysql, postgresql, sqlite3")
	rootCmd.PersistentFlags().StringP("dbURI", "u", "", "Databse URI to connect to, see: http://jinzhu.me/gorm/database.html#connecting-to-a-database")
	viper.BindPFlag("listen", rootCmd.PersistentFlags().Lookup("listen"))
	viper.BindPFlag("dbType", rootCmd.PersistentFlags().Lookup("dbType"))
	viper.BindPFlag("dbURI", rootCmd.PersistentFlags().Lookup("dbURI"))
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc")

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("albiondata-api")

		// Add the executable path as
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		viper.AddConfigPath(exPath)
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
	}

	viper.SetEnvPrefix("ADA")
	viper.AutomaticEnv()
}

func apiHome(c echo.Context) error {
	return c.String(http.StatusOK, "Nothing to show here, please go to https://www.albiononline2d.com/")
}

func apiHandleStatsPricesItem(c echo.Context) error {
	result := []lib.APIStatsPricesItem{}

	itemIDs := strings.Split(c.Param("item"), ",")

	for _, itemID := range itemIDs {
		for _, l := range adslib.Locations() {
			lres := lib.APIStatsPricesItem{
				ItemID: itemID,
				City:   l.String(),
			}

			// Find lowest offer price
			m := adslib.NewModelMarketOrder()
			if err := db.Where("location = ? and item_id = ? and auction_type = ?", l, itemID, "offer").Order("price").First(&m).Error; err != nil {
				continue
			}
			lres.SellPriceMin = m.Price
			lres.SellPriceMinDate = m.UpdatedAt

			// Find highest offer price
			m = adslib.NewModelMarketOrder()
			if err := db.Where("location = ? and item_id = ? and auction_type = ?", l, itemID, "offer").Order("price desc").First(&m).Error; err == nil {
				lres.SellPriceMax = m.Price
				lres.SellPriceMaxDate = m.UpdatedAt
			}

			// Find lowest request price
			m = adslib.NewModelMarketOrder()
			if err := db.Where("location = ? and item_id = ? and auction_type = ?", l, itemID, "request").Order("price").First(&m).Error; err == nil {
				lres.BuyPriceMin = m.Price
				lres.BuyPriceMinDate = m.UpdatedAt
			}

			// Find highest request price
			m = adslib.NewModelMarketOrder()
			if err := db.Where("location = ? and item_id = ? and auction_type = ?", l, itemID, "request").Order("price desc").First(&m).Error; err == nil {
				lres.BuyPriceMax = m.Price
				lres.BuyPriceMaxDate = m.UpdatedAt
			}

			result = append(result, lres)
		}
	}

	return c.JSON(http.StatusOK, result)
}

func doCmd(cmd *cobra.Command, args []string) {
	//******************************
	// START DB
	fmt.Printf("Connecting to database: %s\n", viper.GetString("dbType"))
	var err error
	db, err = gorm.Open(viper.GetString("dbType"), viper.GetString("dbURI"))
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Debug
	// db.LogMode(true)

	defer db.Close()
	// END DB
	//******************************

	//******************************
	// START ECHO
	e := echo.New()
	e.HideBanner = true

	// Recover from panics
	e.Use(middleware.Recover())

	// Logger
	e.Use(middleware.Logger())

	e.GET("/", apiHome)
	e.GET("/api/v1/stats/prices/:item", apiHandleStatsPricesItem)

	// Start server
	e.Logger.Fatal(e.Start(viper.GetString("listen")))

	// END ECHO
	//*******************************
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

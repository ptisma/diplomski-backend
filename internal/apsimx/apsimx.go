package apsimx

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

var Barley = `{
  "$type": "Models.Core.Simulations, Models",
  "ExplorerWidth": 300,
  "Version": 100,
  "ApsimVersion": "0.0.0.0",
  "Name": "Simulations",
  "Children": [
    {
      "$type": "Models.Core.Simulation, Models",
      "IsRunning": false,
      "Name": "Simulation",
      "Children": [
        {
          "$type": "Models.Clock, Models",
          "Start": "%s",
          "End": "%s",
          "Name": "Clock",
          "Children": [],
          "IncludeInDocumentation": true,
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Summary, Models",
          "CaptureErrors": true,
          "CaptureWarnings": true,
          "CaptureSummaryText": true,
          "Name": "SummaryFile",
          "Children": [],
          "IncludeInDocumentation": true,
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Weather, Models",
          "ConstantsFile": "%s",
          "FileName": "%s",
          "ExcelWorkSheetName": "",
          "Name": "Weather",
          "Children": [],
          "IncludeInDocumentation": true,
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Soils.Arbitrator.SoilArbitrator, Models",
          "Name": "SoilArbitrator",
          "Children": [],
          "IncludeInDocumentation": true,
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Core.Zone, Models",
          "Area": 1.0,
          "Slope": 0.0,
          "AspectAngle": 0.0,
          "Altitude": 50.0,
          "Name": "Field",
          "Children": [
            {
              "$type": "Models.Report, Models",
              "VariableNames": [
                "[Clock].Today as Date",
                "[Barley].Grain.Total.Wt*10 as Yield",
                "[Weather].Radn as Radiation",
				"[Weather].MaxT as MaxTemperature",
				"[Weather].MinT as MinTemperature",
				"[Weather].Rain as Rain",
				"[Weather].Wind as Wind",
				"[Weather].Latitude as Latitude",
				"[Weather].Longitude as Longitude",
              ],
              "EventNames": [
                "[Clock].EndOfDay"
              ],
              "GroupByVariableName": null,
              "Name": "Report",
              "Children": [],
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Fertiliser, Models",
              "Name": "Fertiliser",
              "Children": [],
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Soils.Soil, Models",
              "RecordNumber": 0,
              "ASCOrder": "Vertosol",
              "ASCSubOrder": "Black",
              "SoilType": "Clay",
              "LocalName": null,
              "Site": "Norwin",
              "NearestTown": "Norwin",
              "Region": "Darling Downs and Granite Belt",
              "State": "Queensland",
              "Country": "Australia",
              "NaturalVegetation": "Qld. Bluegrass, possible Qld. Blue gum",
              "ApsoilNumber": "900",
              "Latitude": -27.581836,
              "Longitude": 151.320206,
              "LocationAccuracy": " +/- 20m",
              "DataSource": "CSIRO Sustainable Ecosystems, Toowoomba; Characteriesd as part of the GRDC funded project\"Doing it better, doing it smarter, managing soil water in Australian agriculture' 2011",
              "Comments": "OC, CLL for all crops estimated-based on Bongeen Mywybilla Soil No1",
              "Name": "Soil",
              "Children": [
                {
                  "$type": "Models.Soils.Physical, Models",
                  "Depth": [
                    "0-15",
                    "15-30",
                    "30-60",
                    "60-90",
                    "90-120",
                    "120-150",
                    "150-180"
                  ],
                  "Thickness": [
                    150.0,
                    150.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0
                  ],
                  "ParticleSizeClay": null,
                  "ParticleSizeSand": null,
                  "ParticleSizeSilt": null,
                  "BD": [
                    1.01056473311131,
                    1.07145631083388,
                    1.09393858528057,
                    1.15861335018721,
                    1.17301160318016,
                    1.16287303586874,
                    1.18749547755906
                  ],
                  "AirDry": [
                    0.130250054518252,
                    0.198689390775399,
                    0.28,
                    0.28,
                    0.28,
                    0.28,
                    0.28
                  ],
                  "LL15": [
                    0.260500109036505,
                    0.248361738469248,
                    0.28,
                    0.28,
                    0.28,
                    0.28,
                    0.28
                  ],
                  "DUL": [
                    0.52100021807301,
                    0.496723476938497,
                    0.488437607673005,
                    0.480296969355493,
                    0.471583596524955,
                    0.457070570557793,
                    0.452331759845006
                  ],
                  "SAT": [
                    0.588654817693846,
                    0.565676863836273,
                    0.557192986686577,
                    0.532787415023694,
                    0.527354112007486,
                    0.531179986464627,
                    0.521888499034317
                  ],
                  "KS": [
                    20.0,
                    20.0,
                    20.0,
                    20.0,
                    20.0,
                    20.0,
                    20.0
                  ],
                  "BDMetadata": null,
                  "AirDryMetadata": null,
                  "LL15Metadata": null,
                  "DULMetadata": null,
                  "SATMetadata": null,
                  "KSMetadata": null,
                  "Name": "Physical",
                  "Children": [
                    {
                      "$type": "Models.Soils.SoilCrop, Models",
                      "LL": [
                        0.261,
                        0.248,
                        0.28,
                        0.306,
                        0.36,
                        0.392,
                        0.446
                      ],
                      "KL": [
                        0.06,
                        0.06,
                        0.06,
                        0.04,
                        0.04,
                        0.02,
                        0.01
                      ],
                      "XF": [
                        1.0,
                        1.0,
                        1.0,
                        1.0,
                        1.0,
                        1.0,
                        1.0
                      ],
                      "LLMetadata": null,
                      "KLMetadata": null,
                      "XFMetadata": null,
                      "Name": "BarleySoil",
                      "Children": [],
                      "IncludeInDocumentation": true,
                      "Enabled": true,
                      "ReadOnly": false
                    }
                  ],
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.WaterModel.WaterBalance, Models",
                  "SummerDate": "1-Nov",
                  "SummerU": 5.0,
                  "SummerCona": 5.0,
                  "WinterDate": "1-Apr",
                  "WinterU": 5.0,
                  "WinterCona": 5.0,
                  "DiffusConst": 40.0,
                  "DiffusSlope": 16.0,
                  "Salb": 0.12,
                  "CN2Bare": 73.0,
                  "CNRed": 20.0,
                  "CNCov": 0.8,
                  "Slope": "NaN",
                  "DischargeWidth": "NaN",
                  "CatchmentArea": "NaN",
                  "Thickness": [
                    150.0,
                    150.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0
                  ],
                  "SWCON": [
                    0.3,
                    0.3,
                    0.3,
                    0.3,
                    0.3,
                    0.3,
                    0.3
                  ],
                  "KLAT": null,
                  "ResourceName": "WaterBalance",
                  "Name": "SoilWater",
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.Soils.Organic, Models",
                  "Depth": [
                    "0-15",
                    "15-30",
                    "30-60",
                    "60-90",
                    "90-120",
                    "120-150",
                    "150-180"
                  ],
                  "FOMCNRatio": 40.0,
                  "Thickness": [
                    150.0,
                    150.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0
                  ],
                  "Carbon": [
                    1.2,
                    0.96,
                    0.6,
                    0.3,
                    0.18,
                    0.12,
                    0.12
                  ],
                  "SoilCNRatio": [
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0,
                    12.0
                  ],
                  "FBiom": [
                    0.04,
                    0.02,
                    0.02,
                    0.02,
                    0.01,
                    0.01,
                    0.01
                  ],
                  "FInert": [
                    0.4,
                    0.6,
                    0.8,
                    1.0,
                    1.0,
                    1.0,
                    1.0
                  ],
                  "FOM": [
                    347.12903231275641,
                    270.3443621919937,
                    163.97214434990104,
                    99.454132887040629,
                    60.321980831124677,
                    36.587130828674873,
                    22.1912165985086
                  ],
                  "Name": "Organic",
                  "Children": [],
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.Soils.Chemical, Models",
                  "Depth": [
                    "0-15",
                    "15-30",
                    "30-60",
                    "60-90",
                    "90-120",
                    "120-150",
                    "150-180"
                  ],
                  "Thickness": [
                    150.0,
                    150.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0
                  ],
                  "NO3N": [
                    1.0,
                    1.0,
                    1.0,
                    1.0,
                    1.0,
                    1.0,
                    1.0
                  ],
                  "NH4N": [
                    0.1,
                    0.1,
                    0.1,
                    0.1,
                    0.1,
                    0.1,
                    0.1
                  ],
                  "PH": [
                    8.0,
                    8.0,
                    8.0,
                    8.0,
                    8.0,
                    8.0,
                    8.0
                  ],
                  "CL": null,
                  "EC": null,
                  "ESP": null,
                  "Name": "Chemical",
                  "Children": [],
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.Soils.InitialWater, Models",
                  "PercentMethod": 1,
                  "FractionFull": 1.0,
                  "DepthWetSoil": "NaN",
                  "RelativeTo": null,
                  "Name": "InitialWater",
                  "Children": [],
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.Soils.Sample, Models",
                  "Depth": [
                    "0-15",
                    "15-30",
                    "30-60",
                    "60-90",
                    "90-120",
                    "120-150",
                    "150-180"
                  ],
                  "Thickness": [
                    150.0,
                    150.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0,
                    300.0
                  ],
                  "NO3N": null,
                  "NH4N": null,
                  "SW": [
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN"
                  ],
                  "OC": [
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN"
                  ],
                  "EC": [
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN"
                  ],
                  "CL": [
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN"
                  ],
                  "ESP": [
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN"
                  ],
                  "PH": [
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN",
                    "NaN"
                  ],
                  "SWUnits": 0,
                  "OCUnits": 0,
                  "PHUnits": 0,
                  "Name": "InitialN",
                  "Children": [],
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.Soils.CERESSoilTemperature, Models",
                  "Name": "CERESSoilTemperature",
                  "Children": [],
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                },
                {
                  "$type": "Models.Soils.Nutrients.Nutrient, Models",
                  "ResourceName": "Nutrient",
                  "Name": "Nutrient",
                  "IncludeInDocumentation": true,
                  "Enabled": true,
                  "ReadOnly": false
                }
              ],
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Surface.SurfaceOrganicMatter, Models",
              "InitialResidueName": "wheat_stubble",
              "InitialResidueType": "wheat",
              "InitialResidueMass": 500.0,
              "InitialStandingFraction": 0.0,
              "InitialCPR": 0.0,
              "InitialCNR": 100.0,
              "ResourceName": "SurfaceOrganicMatter",
              "Name": "SurfaceOrganicMatter",
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.MicroClimate, Models",
              "a_interception": 0.0,
              "b_interception": 1.0,
              "c_interception": 0.0,
              "d_interception": 0.0,
              "soil_albedo": 0.3,
              "SoilHeatFluxFraction": 0.4,
              "MinimumHeightDiffForNewLayer": 0.0,
              "NightInterceptionFraction": 0.5,
              "ReferenceHeight": 2.0,
              "Name": "MicroClimate",
              "Children": [],
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Manager, Models",
              "Code": "using Models.PMF;\r\nusing Models.Core;\r\nusing System;\r\nnamespace Models\r\n{\r\n    [Serializable]\r\n    public class Script : Model\r\n    {\r\n        [Link] Clock Clock;\r\n        [Link] Fertiliser Fertiliser;\r\n        [Link] Summary Summary;\r\n        \r\n        \r\n        [Description(\"Amount of fertiliser to be applied (kg/ha)\")]\r\n        public double Amount { get; set;}\r\n        \r\n        [Description(\"Crop to be fertilised\")]\r\n        public string CropName { get; set;}\r\n        \r\n        \r\n        \r\n\r\n        [EventSubscribe(\"Sowing\")]\r\n        private void OnSowing(object sender, EventArgs e)\r\n        {\r\n            Model crop = sender as Model;\r\n            if (crop.Name.ToLower()==CropName.ToLower())\r\n                Fertiliser.Apply(Amount: Amount, Type: Fertiliser.Types.NO3N);\r\n        }\r\n        \r\n    }\r\n}\r\n",
              "Parameters": [
                {
                  "Key": "Amount",
                  "Value": "160"
                },
                {
                  "Key": "CropName",
                  "Value": "Barley"
                }
              ],
              "Name": "SowingFertiliser",
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Manager, Models",
              "Code": "using APSIM.Shared.Utilities;\r\nusing Models.Utilities;\r\nusing Models.Soils.Nutrients;\r\nusing Models.Soils;\r\nusing Models.PMF;\r\nusing Models.Core;\r\nusing System;\r\n\r\nnamespace Models\r\n{\r\n    [Serializable]\r\n    public class Script : Model\r\n    {\r\n        [Link(ByName = true)] Plant Barley;\r\n\r\n        [EventSubscribe(\"DoManagement\")]\r\n        private void OnDoManagement(object sender, EventArgs e)\r\n        {\r\n            if (Barley.IsReadyForHarvesting)\r\n            {\r\n               Barley.Harvest();\r\n               Barley.EndCrop();    \r\n            }\r\n        \r\n        }\r\n        \r\n    }\r\n}\r\n",
              "Parameters": [],
              "Name": "Harvest",
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.PMF.Plant, Models",
              "ResourceName": "Barley",
              "Name": "Barley",
              "IncludeInDocumentation": false,
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Manager, Models",
              "Code": "using APSIM.Shared.Utilities;\r\nusing Models.Utilities;\r\nusing Models.Soils.Nutrients;\r\nusing Models.Soils;\r\nusing Models.PMF;\r\nusing Models.Core;\r\nusing System;\r\n\r\nnamespace Models\r\n{\r\n    [Serializable]\r\n    public class Script : Model\r\n    {\r\n        [Link] Clock Clock;\r\n        [Link] Fertiliser Fertiliser;\r\n        [Link] Summary Summary;\r\n        [Link] Soil Soil;\r\n        Accumulator accumulatedRain;\r\n        \r\n        [Description(\"Crop\")]\r\n        public IPlant Crop { get; set; }\r\n        [Description(\"Sowing date (d-mmm)\")]\r\n        public string SowDate { get; set; }\r\n    [Display(Type = DisplayType.CultivarName, PlantName = \"Barley\")]\r\n        [Description(\"Cultivar to be sown\")]\r\n        public string CultivarName { get; set; }\r\n        [Description(\"Sowing depth (mm)\")]\r\n        public double SowingDepth { get; set; }\r\n        [Description(\"Row spacing (mm)\")]\r\n        public double RowSpacing { get; set; }\r\n        [Description(\"Plant population (/m2)\")]\r\n        public double Population { get; set; }\r\n        \r\n\r\n\r\n        [EventSubscribe(\"DoManagement\")]\r\n        private void OnDoManagement(object sender, EventArgs e)\r\n        {\r\n            if (DateUtilities.WithinDates(SowDate, Clock.Today, SowDate))\r\n            {\r\n                Crop.Sow(population: Population, cultivar: CultivarName, depth: SowingDepth, rowSpacing: RowSpacing);    \r\n            }\r\n        \r\n        }\r\n        \r\n    }\r\n}\r\n",
              "Parameters": [
                {
                  "Key": "Crop",
                  "Value": "[Barley]"
                },
                {
                  "Key": "SowDate",
                  "Value": "24-Jul"
                },
                {
                  "Key": "CultivarName",
                  "Value": "Dash"
                },
                {
                  "Key": "SowingDepth",
                  "Value": "50"
                },
                {
                  "Key": "RowSpacing",
                  "Value": "750"
                },
                {
                  "Key": "Population",
                  "Value": "6"
                }
              ],
              "Name": "Sow on a fixed date",
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            }
          ],
          "IncludeInDocumentation": true,
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Graph, Models",
          "Caption": null,
          "Axis": [
            {
              "$type": "Models.Axis, Models",
              "Type": 3,
              "Title": "Date",
              "Inverted": false,
              "Minimum": "NaN",
              "Maximum": "NaN",
              "Interval": "NaN",
              "DateTimeAxis": false,
              "CrossesAtZero": false
            },
            {
              "$type": "Models.Axis, Models",
              "Type": 0,
              "Title": "Yield (kg/ha)",
              "Inverted": false,
              "Minimum": "NaN",
              "Maximum": "NaN",
              "Interval": "NaN",
              "DateTimeAxis": false,
              "CrossesAtZero": false
            }
          ],
          "LegendPosition": 0,
          "LegendOrientation": 0,
          "DisabledSeries": [],
          "LegendOutsideGraph": false,
          "Name": "Barley Yield Time Series",
          "Children": [
            {
              "$type": "Models.Series, Models",
              "Type": 1,
              "XAxis": 3,
              "YAxis": 0,
              "ColourArgb": -16777216,
              "FactorToVaryColours": null,
              "FactorToVaryMarkers": null,
              "FactorToVaryLines": null,
              "Marker": 0,
              "MarkerSize": 0,
              "Line": 0,
              "LineThickness": 0,
              "TableName": "Report",
              "XFieldName": "Clock.Today",
              "YFieldName": "Yield",
              "X2FieldName": "",
              "Y2FieldName": "",
              "ShowInLegend": true,
              "IncludeSeriesNameInLegend": false,
              "Cumulative": false,
              "CumulativeX": false,
              "Filter": null,
              "Name": "Barley Yield",
              "Children": [],
              "IncludeInDocumentation": true,
              "Enabled": true,
              "ReadOnly": false
            }
          ],
          "IncludeInDocumentation": true,
          "Enabled": true,
          "ReadOnly": false
        }
      ],
      "IncludeInDocumentation": true,
      "Enabled": true,
      "ReadOnly": false
    },
    {
      "$type": "Models.Storage.DataStore, Models",
      "useFirebird": false,
      "CustomFileName": null,
      "Name": "DataStore",
      "Children": [],
      "IncludeInDocumentation": true,
      "Enabled": true,
      "ReadOnly": false
    }
  ],
  "IncludeInDocumentation": true,
  "Enabled": true,
  "ReadOnly": false
}`

//API
func NewBarleyApsimx(start, end time.Time, weatherFilePath string) string {

	if runtime.GOOS == "windows" {
		weatherFilePath = strings.Replace(weatherFilePath, "\\", "\\\\", -1)
	}
	fmt.Println("weatherFilePath:", weatherFilePath)

	//same format as apsimx example file
	res := fmt.Sprintf(Barley, start.Format("2006-01-02T15:04:05"), end.Format("2006-01-02T15:04:05"), weatherFilePath)
	return res

}

//CSV
func NewBarleyApsimx2(start, end time.Time, csvFilePath, constFilePath string) string {

	if runtime.GOOS == "windows" {
		csvFilePath = strings.Replace(csvFilePath, "\\", "\\\\", -1)
		constFilePath = strings.Replace(constFilePath, "\\", "\\\\", -1)

	}
	fmt.Println("csvFilePath:", csvFilePath)
	fmt.Println("constFilePath:", constFilePath)

	//same format as apsimx example file
	res := fmt.Sprintf(Barley, start.Format("2006-01-02T15:04:05"), end.Format("2006-01-02T15:04:05"), constFilePath, csvFilePath)
	return res

}

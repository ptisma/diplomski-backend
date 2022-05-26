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
			%s,
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
                  "Value": "24-Nov"
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

//TODO
var Maize = `{
  "$type": "Models.Core.Simulations, Models",
  "ExplorerWidth": 300,
  "Version": 148,
  "Name": "Simulations",
  "Children": [
    {
      "$type": "Models.Core.Simulation, Models",
      "IsRunning": false,
      "Descriptors": null,
      "Name": "Simulation",
      "Children": [
        {
          "$type": "Models.Clock, Models",
          "Start": "%s",
          "End": "%s",
          "Name": "Clock",
          "Children": [],
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Summary, Models",
          "Verbosity": 100,
          "Name": "SummaryFile",
          "Children": [],
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Climate.Weather, Models",
          "ConstantsFile": "%s",
          "FileName": "%s",
          "ExcelWorkSheetName": "",
          "Name": "Weather",
          "Children": [],
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Soils.Arbitrator.SoilArbitrator, Models",
          "Name": "Soil Arbitrator",
          "Children": [],
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
                "[Clock].Today as Date,",
                "[Maize].Grain.Total.Wt*10 as Yield,",
                "[Weather].Radn as Radiation,",
                "[Weather].MaxT as MaxTemperature,",
                "[Weather].MinT as MinTemperature,",
                "[Weather].Rain as Rain,",
                "[Weather].Wind as Wind,",
                "[Weather].Latitude as Latitude,",
                "[Weather].Longitude as Longitude,"
              ],
              "EventNames": [
                "[Clock].EndOfDay"
              ],
              "GroupByVariableName": null,
              "Name": "Report",
              "Children": [],
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Fertiliser, Models",
              "ResourceName": "Fertiliser",
              "Name": "Fertiliser",
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
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.PMF.Plant, Models",
              "ResourceName": "Maize",
              "Name": "Maize",
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.MicroClimate, Models",
              "a_interception": 0.0,
              "b_interception": 1.0,
              "c_interception": 0.0,
              "d_interception": 0.0,
              "SoilHeatFluxFraction": 0.4,
              "MinimumHeightDiffForNewLayer": 0.0,
              "NightInterceptionFraction": 0.5,
              "ReferenceHeight": 2.0,
              "Name": "MicroClimate",
              "Children": [],
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Manager, Models",
              "Code": "using Models.PMF;\r\nusing Models.Core;\r\nusing System;\r\nnamespace Models\r\n{\r\n    [Serializable]\r\n    public class Script : Model\r\n    {\r\n        [Link] Clock Clock;\r\n        [Link] Fertiliser Fertiliser;\r\n        [Link] Summary Summary;\r\n        \r\n        \r\n        [Description(\"Amount of fertiliser to be applied (kg/ha)\")]\r\n        public double Amount { get; set; }\r\n        \r\n        [Description(\"Crop to be fertilised\")]\r\n        public string CropName { get; set; }\r\n        \r\n        \r\n        \r\n\r\n        [EventSubscribe(\"Sowing\")]\r\n        private void OnSowing(object sender, EventArgs e)\r\n        {\r\n            Model crop = sender as Model;\r\n            if (crop.Name.ToLower() == CropName.ToLower())\r\n                Fertiliser.Apply(Amount: Amount, Type: Fertiliser.Types.NO3N);\r\n        }\r\n        \r\n    }\r\n}\r\n",
              "Parameters": [
                {
                  "Key": "Amount",
                  "Value": "160"
                },
                {
                  "Key": "CropName",
                  "Value": "maize"
                }
              ],
              "Name": "SowingFertiliser",
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Manager, Models",
              "Code": "using Models.Interfaces;\r\nusing System;\r\nusing Models.Core;\r\nusing Models.PMF;\r\nusing Models.Soils;\r\nusing Models.Soils.Nutrients;\r\nusing Models.Utilities;\r\nusing APSIM.Shared.Utilities;\r\nusing Models.Climate;\r\n\r\nnamespace Models\r\n{\r\n    [Serializable]\r\n    public class Script : Model\r\n    {\r\n        [Link] private Clock Clock;\r\n        [Link] private Fertiliser Fertiliser;\r\n        [Link] private Summary Summary;\r\n        [Link(ByName = true)] private Plant Maize;\r\n        [Link] private Soil Soil;\r\n        private Accumulator accumulatedRain;\r\n        [Link]\r\n        private ISoilWater waterBalance;\r\n        \r\n        [Description(\"Start of sowing window (d-mmm)\")]\r\n        public string StartDate { get; set; }\r\n        [Description(\"End of sowing window (d-mmm)\")]\r\n        public string EndDate { get; set; }\r\n        [Description(\"Minimum extractable soil water for sowing (mm)\")]\r\n        public double MinESW { get; set; }\r\n        [Description(\"Accumulated rainfall required for sowing (mm)\")]\r\n        public double MinRain { get; set; }\r\n        [Description(\"Duration of rainfall accumulation (d)\")]\r\n        public int RainDays { get; set; }\r\n        [Description(\"Cultivar to be sown\")]\r\n        [Display(Type=DisplayType.CultivarName, PlantName = \"Maize\")]\r\n        public string CultivarName { get; set; }\r\n        [Description(\"Sowing depth (mm)\")]\r\n        public double SowingDepth { get; set; }\r\n        [Description(\"Row spacing (mm)\")]\r\n        public double RowSpacing { get; set; }\r\n        [Description(\"Plant population (/m2)\")]\r\n        public double Population { get; set; }\r\n        \r\n        \r\n        [EventSubscribe(\"StartOfSimulation\")]\r\n        private void OnSimulationCommencing(object sender, EventArgs e)\r\n        {\r\n            accumulatedRain = new Accumulator(this, \"[Weather].Rain\", RainDays);\r\n        }\r\n        \r\n\r\n        [EventSubscribe(\"DoManagement\")]\r\n        private void OnDoManagement(object sender, EventArgs e)\r\n        {\r\n            accumulatedRain.Update();\r\n            \r\n            if (DateUtilities.WithinDates(StartDate, Clock.Today, EndDate) &&\r\n                !Maize.IsAlive &&\r\n                MathUtilities.Sum(waterBalance.ESW) > MinESW &&\r\n                accumulatedRain.Sum > MinRain)\r\n            {\r\n                Maize.Sow(population: Population, cultivar: CultivarName, depth: SowingDepth, rowSpacing: RowSpacing);    \r\n            }\r\n        \r\n        }\r\n        \r\n    }\r\n}\r\n",
              "Parameters": [
                {
                  "Key": "StartDate",
                  "Value": "15-apr"
                },
                {
                  "Key": "EndDate",
                  "Value": "25-may"
                },
                {
                  "Key": "MinESW",
                  "Value": "100"
                },
                {
                  "Key": "MinRain",
                  "Value": "25"
                },
                {
                  "Key": "RainDays",
                  "Value": "7"
                },
                {
                  "Key": "CultivarName",
                  "Value": "Dekalb_XL82"
                },
                {
                  "Key": "SowingDepth",
                  "Value": "30"
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
              "Name": "SowingRule",
              "Enabled": true,
              "ReadOnly": false
            },
            {
              "$type": "Models.Manager, Models",
              "Code": "using Models.PMF.Phen;\r\nusing APSIM.Shared.Utilities;\r\nusing System.Xml.Serialization;\r\nusing Models;\r\nusing Models.PMF;\r\nusing Models.Soils;\r\nusing System.Text;\r\nusing System.Collections.Generic;\r\nusing Models.Core;\r\nusing System;\r\nusing Models.Soils.Nutrients;\r\n\r\nnamespace Models\r\n{\r\n    [Serializable] \r\n    public class Script : Model\r\n    {\r\n        [Link(ByName = true)] private Plant Maize;\r\n        [Link] private Zone zone;\r\n        [Link(Type = LinkType.Path, Path = \"[Maize].Phenology\")]\r\n        private Phenology phenology;\r\n                    \r\n        \r\n        [EventSubscribe(\"StartOfSimulation\")]\r\n        private void OnSimulationCommencing(object sender, EventArgs e)\r\n        {\r\n        }\r\n        \r\n        [EventSubscribe(\"DoManagement\")]\r\n        private void OnDoManagement(object sender, EventArgs e)\r\n        {\r\n            if (phenology.CurrentPhase.Name == \"ReadyForHarvesting\")\r\n            {\r\n                Maize.Harvest();\r\n                Maize.EndCrop();\r\n            }\r\n        }\r\n        \r\n        [EventSubscribe(\"DoManagementCalculations\")]\r\n        private void OnDoManagementCalculations(object sender, EventArgs e)\r\n        {\r\n            \r\n        }\r\n    }\r\n}\r\n       \r\n",
              "Parameters": [],
              "Name": "Harvesting",
              "Enabled": true,
              "ReadOnly": false
            },
            %s
          ],
          "Enabled": true,
          "ReadOnly": false
        },
        {
          "$type": "Models.Graph, Models",
          "Caption": null,
          "Axis": [
            {
              "$type": "APSIM.Shared.Graphing.Axis, APSIM.Shared",
              "Title": "Date",
              "Position": 3,
              "Inverted": false,
              "CrossesAtZero": false,
              "Minimum": null,
              "Maximum": null,
              "Interval": null
            },
            {
              "$type": "APSIM.Shared.Graphing.Axis, APSIM.Shared",
              "Title": null,
              "Position": 0,
              "Inverted": false,
              "CrossesAtZero": false,
              "Minimum": null,
              "Maximum": null,
              "Interval": null
            }
          ],
          "LegendPosition": 0,
          "LegendOrientation": 0,
          "AnnotationLocation": 0,
          "DisabledSeries": [],
          "LegendOutsideGraph": false,
          "Name": "Maize Yield Time Series",
          "Children": [
            {
              "$type": "Models.Series, Models",
              "Type": 1,
              "XAxis": 3,
              "YAxis": 0,
              "ColourArgb": -16776961,
              "FactorToVaryColours": null,
              "FactorToVaryMarkers": null,
              "FactorToVaryLines": null,
              "Marker": 0,
              "MarkerSize": 0,
              "Line": 0,
              "LineThickness": 0,
              "TableName": "Report",
              "XFieldName": "Clock.Today",
              "YFieldName": "Maize.Grain.Wt",
              "X2FieldName": "",
              "Y2FieldName": "",
              "ShowInLegend": true,
              "IncludeSeriesNameInLegend": false,
              "Cumulative": false,
              "CumulativeX": false,
              "Filter": null,
              "Name": "Maize Yield",
              "Children": [],
              "Enabled": true,
              "ReadOnly": false
            }
          ],
          "Enabled": true,
          "ReadOnly": false
        }
      ],
      "Enabled": true,
      "ReadOnly": false
    },
    {
      "$type": "Models.Storage.DataStore, Models",
      "useFirebird": false,
      "CustomFileName": null,
      "Name": "DataStore",
      "Children": [],
      "Enabled": true,
      "ReadOnly": false
    }
  ],
  "Enabled": true,
  "ReadOnly": false
}`

type GenerateAPSIMCulture func(start, end time.Time, csvFilePath, constFilePath string, soilData string) string

var InitiateAPSIMCulture = map[int]GenerateAPSIMCulture{
	1: NewBarleyApsimx2,
	2: NewMaizeApsimx,
}

//API
//func NewBarleyApsimx(start, end time.Time, weatherFilePath string) string {
//
//	if runtime.GOOS == "windows" {
//		weatherFilePath = strings.Replace(weatherFilePath, "\\", "\\\\", -1)
//	}
//	fmt.Println("weatherFilePath:", weatherFilePath)
//
//	//same format as apsimx example file
//	res := fmt.Sprintf(Barley, start.Format("2006-01-02T15:04:05"), end.Format("2006-01-02T15:04:05"), weatherFilePath)
//	return res
//
//}

//CSV
func NewBarleyApsimx2(start, end time.Time, csvFilePath, constFilePath string, soilData string) string {

	if runtime.GOOS == "windows" {
		csvFilePath = strings.Replace(csvFilePath, "\\", "\\\\", -1)
		constFilePath = strings.Replace(constFilePath, "\\", "\\\\", -1)

	}
	fmt.Println("csvFilePath:", csvFilePath)
	fmt.Println("constFilePath:", constFilePath)

	//same format as apsimx example file
	res := fmt.Sprintf(Barley, start.Format("2006-01-02T15:04:05"), end.Format("2006-01-02T15:04:05"), constFilePath, csvFilePath, soilData)
	return res

}

func NewMaizeApsimx(start, end time.Time, csvFilePath, constFilePath string, soilData string) string {

	if runtime.GOOS == "windows" {
		csvFilePath = strings.Replace(csvFilePath, "\\", "\\\\", -1)
		constFilePath = strings.Replace(constFilePath, "\\", "\\\\", -1)

	}
	fmt.Println("csvFilePath:", csvFilePath)
	fmt.Println("constFilePath:", constFilePath)

	//same format as apsimx example file
	res := fmt.Sprintf(Maize, start.Format("2006-01-02T15:04:05"), end.Format("2006-01-02T15:04:05"), constFilePath, csvFilePath, soilData)
	return res

}

# Instructions

Hi, curious folks! I bet you are one of those people, who is trying very hard to make
this world a better place by providing better health care support to people. Follow
the instructions to setup a local Open Health Algorithms Service into your server and
use it to provide advance health care support.

## Contents in the Zip File

When you unzip the `ohas.zip` into your favorite directory, you will see these files:

```text
|-- ohas
    |-- guideline_hearts.json
    |   (Has the conditions/targets for various health check)
    |-- guideline_hearts_content.json
    |   (Has the care plan messages for different attributes)
    |-- INSTRUCTIONS.md
    |   (this file)
    |-- ohas-darwin-amd64.bin
    |   (The binary you can use in Mac)
    |-- ohas-linux-386.bin
    |   (The binary you can use in 32 bit Linux)
    |-- ohas-linux-amd64.bin
    |   (The binary you can use in 64 bit Linux)
    |-- sample-request.json
    |   (A sample request object)
```

## How to Run

Assuming you unzipped the contents in `~/ohas` directory, then

### Start Service

Run the following command:

```bash
cd ~/ohas && ./ohas-darwin-amd64.bin start
```

### Stop Service

Run the following command:

```bash
cd ~/ohas && ./ohas-darwin-amd64.bin stop
```

## API Usage

The service will run on port `9595` or your given port on the server you installed.
You can use postman or any application to use the API.

### `/api/algorithm`

This is the primary API to run the algorithm.

Here is a sample request:

```bash
curl -X POST \
  http://localhost:9595/api/algorithm \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8a411071-c296-4d65-ba0d-11b3b6425b27' \
  -d '{
  "region": "SEARD",
  "demographics": {
    "gender": "F",
    "age": 60,
    "dob": [
      "computed",
      "01/10/1987"
    ]
  },
  "measurements": {
    "height": [
      1.5,
      "m"
    ],
    "weight": [
      70,
      "kg"
    ],
    "waist": [
      99,
      "cm"
    ],
    "hip": [
      104,
      "cm"
    ],
    "sbp": [
      100,
      120,
      130,
      "sitting"
    ],
    "dbp": [
      91,
      95,
      85,
      "sitting"
    ]
  },
  "smoking": {
    "current": 0,
    "ex_smoker": 1,
    "quit_within_year": 0
  },
  "physical_activity": "120",
  "diet_history": {
    "fruit": 1,
    "veg": 6,
    "rice": 2,
    "oil": "olive"
  },
  "medical_history": {
    "conditions": [
      "asthma",
      "tuberculosis"
    ]
  },
  "allergies": {},
  "medications": [
    "anti_hypertensive",
    "statin",
    "antiplatelet",
    "bronchodilator"
  ],
  "family_history": [
    "cvd"
  ],
  "pathology": {
    "bsl": {
      "type": "random",
      "units": "mg/dl",
      "value": 180
    },
    "cholesterol": {
      "type": "fasting",
      "units": "mmol/l",
      "total_chol": 5.2,
      "hdl": 100,
      "ldl": 240
    }
  }
}
'
```

This is the returned output:

```javascript
{
    "hearts": {
        "blood_pressure": {
            "bp": "116/90",
            "code": "BP-3A",
            "output": [
                "on_target_dm",
                "On target for patient with diabetes",
                "GREEN",
                "Good job"
            ],
            "target": "130/80"
        },
        "cvd_assessment": {
            "cvd_risk_result": {
                "risk": "20",
                "risk_range": "10-20"
            },
            "guidelines": {
                "advice": [
                    "lifestyle",
                    "review-targets"
                ],
                "bp_target": "",
                "follow_up_interval": "3",
                "follow_up_message": "",
                "label": "Low-Moderate",
                "management": {
                    "lifestyle": "Lifestyle modification is an important component of disease management. Focus on healthy diet and regula rexercise. Discuss risk factors",
                    "review-targets": "Keep an eye on your numbers and make sure to review with your doctor."
                },
                "score": "10-20%"
            },
            "high_risk_condition": {
                "code": "",
                "reason": "",
                "status": false
            }
        },
        "diabetes": {
            "code": "DM-3",
            "output": [
                "new_diagnosis",
                "New Diagnosis of DM",
                "RED",
                "Looks like newly diagnosed diabetes"
            ],
            "status": true,
            "value": 10
        },
        "lifestyle": {
            "bmi": {
                "code": "BMI-3",
                "output": [
                    "obese",
                    "Weight in the obese range",
                    "RED",
                    "Obesity is a risk factor for developing many conditions. Talk to a dietitian today"
                ],
                "target": "18.5 - 24.9",
                "value": "31.11"
            },
            "diet": {
                "code": "NUT-2",
                "output": [
                    "fv_off_target_partial",
                    "Partially meeting targets",
                    "YELLOW",
                    "Almost there .. Try to meet the recommended targets for both Fruit and Vegetables"
                ],
                "value": {
                    "fruit": 1,
                    "vegetables": 6
                }
            },
            "exercise": {
                "code": "PA-2",
                "output": [
                    "off_target_mild",
                    "Off target, but trying",
                    "AMBER",
                    "Keep moving. Remember every bit counts. Aim for 150 minutes of moderate intensity exercsie per week. [Learn more]"
                ],
                "target": "150 minutes",
                "value": 120
            },
            "smoking": {
                "code": "SM-3",
                "output": [
                    "ex_smoker",
                    "Ex smoker, quit > 12 months ago",
                    "GREEN",
                    "Good work quitting. It\\\"s the best thing you can do for your health. Stay strong"
                ],
                "smoking_calc": false,
                "status": false
            },
            "whr": {
                "code": "WHR-1",
                "output": [
                    "abnormal_female",
                    "Abnormal waist hip ratio for female",
                    "AMBER",
                    "Your wait hip ratio is above the target for females"
                ],
                "target": "0.85",
                "value": "0.95"
            }
        }
    },
    "request_id": "47054119-d21a-4842-b1bf-ce9cbfd19eb7"
}
```

### `/api/version`

It will return the API version.

```javascript
{
    "version": "0.1"
}
```

## FAQ

- The binary file and both guideline files must be on the same directory.
- Do not rename the guideline files.
- Most of the attributes in `sample-request.json` are mandatory. If you get invalid/unexpected output, the check the input attributes
- Many error checking is not fully implemented, so use in your own risk.

## Coming Soon

- Configurable Service
- Loggers
- Service Installer
- FHIR Compatibility

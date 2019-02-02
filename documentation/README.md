# Open Health Algorithms Service

This service will generate binaries to be used as a package in any posix system. It will expose a `service` in `9595` port.

## Table of Contents

- [Installation](#installation)
- [API End Points](#api-end-points)
- [Algorithm Input Parameters](#algorithm-input-parameters)
- [Algorithm Output](#algorithm-output)
- [Errors](#errors)

## Installation

Installing the service is very simple as it does not need to run any installer. After you got the latest version of the service
in a zip file, unzip to a directory of your choice. This zip file contains the following files and folders inside:

```text
|-- ohas
    |-- documentation
        |-- assets/
        |-- icon.png
        |-- index.html
        |-- README.md (This file)
    |-- goals_hearts_content.json (Has the goals content for various health check)
    |-- goals_hearts.json (Has the goals conditions for various health check)
    |-- guideline_hearts_content.json (Has the care plan messages for different attributes)
    |-- guideline_hearts.json (Has the conditions/targets for various health check)
    |-- ohas-darwin-amd64.bin (The binary you can use in Mac)
    |-- ohas-linux-amd64.bin (The binary you can use in Linux)
    |-- sample-request.json (A sample request object)
```

Assuming you unzipped the contents in `~/ohas` directory, then to start the service, run the following command:

```bash
cd ~/ohas && ./ohas-linux-amd64.bin start
```

To stop Service, run the following command:

```bash
cd ~/ohas && ./ohas-linux-amd64.bin stop
```

To view the documentation in the given html file, run it within a webserver.
For example, you can run a simple server for preview:

```bash
php -S localhost:7654 -t ~/ohas/documentation
```

Then view [http://localhost:7654](http://localhost:7654)

## API End Points

The service has two API end points at the moment.

### Version API `GET /api/version`

This API will return the current version of the service. It uses Semantic Versioning (SemVer) format.

```bash
curl -X GET https://demoservice.openhealthalgorithms.org/api/version
```

The **response** will be a json object with an attribute `version`.

```javascript
{
    "version": "0.4.3"
}
```

### Algorithm API `POST /api/algorithm`

This is the primary API to run the algorithm.

Here is a sample request:

```bash
curl -X POST \
  https://demoservice.openhealthalgorithms.org/api/algorithm \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
  "config": {
    "algorithm": "hearts",
    "risk_model": "whocvd"
  },
  "params": {
    "demographics": {
      "gender": "male",
      "age": {
        "value": 65,
        "unit": "year"
      },
      "birth_country": "Bangladesh",
      "birth_country_code": "BD",
      "living_country": "Bangladesh",
      "living_country_code": "BD",
      "race": "Bengali",
      "ethnicity": "Bengali"
    },
    "components": {
      "lifestyle": [
        {
          "name": "smoking",
          "category": "addiction",
          "value": "ex-smoker",
          "quit_within_year": true
        },
        {
          "name": "alcohol_history",
          "category": "addiction",
          "value": 12,
          "units": "units",
          "frequency": "weekly"
        },
        {
          "name": "fruit",
          "category": "nutrition",
          "value": 3,
          "units": "servings",
          "frequency": "daily"
        },
        {
          "name": "vegetables",
          "category": "nutrition",
          "value": 2,
          "units": "servings",
          "frequency": "daily"
        },
        {
          "name": "low-intensity-exercise",
          "category": "physical-activity",
          "value": 20,
          "units": "minutes",
          "frequency": "daily"
        }
      ],
      "body-measurements": [
        {
          "effectiveDate": "timestamp",
          "name": "blood_pressure",
          "category": "vital-signs",
          "value": "170/80",
          "units": "mmHg",
          "arm": "right"
        },
        {
          "effectiveDate": "timestamp",
          "name": "blood_pressure",
          "category": "vital-signs",
          "value": "130/85",
          "units": "mmHg",
          "arm": "right"
        },
        {
          "effectiveDate": "timestamp",
          "name": "blood_pressure",
          "category": "vital-signs",
          "value": "135/90",
          "units": "mmHg",
          "arm": "right"
        },
        {
          "effectiveDate": "timestamp",
          "name": "height",
          "category": "anthropometry",
          "units": "cm",
          "value": 180.5
        },
        {
          "effectiveDate": "timestamp",
          "name": "weight",
          "category": "anthropometry",
          "units": "kg",
          "value": 80.5
        },
        {
          "effectiveDate": "timestamp",
          "name": "hip",
          "category": "anthropometry",
          "units": "cm",
          "value": 110
        },
        {
          "effectiveDate": "timestamp",
          "name": "waist",
          "category": "anthropometry",
          "units": "cm",
          "value": 100.7
        },
        {
          "effectiveDate": "timestamp",
          "name": "body-fat",
          "category": "anthropometry",
          "units": "%",
          "value": 25
        },
        {
          "effectiveDate": "timestamp",
          "name": "muscle",
          "category": "anthropometry",
          "units": "%",
          "value": 20
        },
        {
          "effectiveDate": "timestamp",
          "name": "visceral_fat",
          "category": "anthropometry",
          "units": "%",
          "value": 10
        }
      ],
      "biological-samples": [
        {
          "effectiveDate": "timestamp",
          "name": "blood_sugar",
          "category": "blood-test",
          "value": 15.6,
          "units": "mmol/L",
          "type": "fasting"
        },
        {
          "effectiveDate": "timestamp",
          "name": "total_cholesterol",
          "category": "blood-test",
          "value": 5.6,
          "units": "mmol/L"
        },
        {
          "effectiveDate": "timestamp",
          "name": "hdl",
          "category": "blood-test",
          "value": 1.3,
          "units": "mmol/L"
        },
        {
          "effectiveDate": "timestamp",
          "name": "ldl",
          "category": "blood-test",
          "value": 1.3,
          "units": "mmol/L"
        },
        {
          "effectiveDate": "timestamp",
          "name": "tg",
          "category": "blood-test",
          "value": 6,
          "units": "mmol/L"
        }
      ],
      "medical_history": [
        {
          "name": "hypertension",
          "category": "condition",
          "is_active": true
        },
        {
          "name": "diabetes",
          "category": "condition",
          "is_active": true
        },
        {
          "name": "tuberculosis",
          "category": "condition",
          "is_active": false
        },
        {
          "category": "allergy",
          "type": "medication",
          "allergen": "ace-i",
          "criticality": "",
          "reaction": "cough"
        }
      ],
      "medications": [
        {
          "generic": "rampiril",
          "category": "anti-hypertensive"
        }
      ],
      "family_history": [
        {
          "name": "cardiovascular-disease",
          "relative": "1st degree"
        }
      ]
    }
  }
}
'
```

The **response** will be a json object.

```javascript
{
    "errors": [],
    "hearts": {
        "assessments": {
            "blood_pressure": {
                "code": "BP-HTN",
                "eval": "Hypertension",
                "grading": 0,
                "message": "Your BP indicates hypertension. Discuss with local health professional",
                "refer": "yes",
                "target": "130/80",
                "tfl": "AMBER",
                "value": "145/85"
            },
            "body_composition": {
                "components": {
                    "bmi": {
                        "code": "BMI-OVERWEIGHT",
                        "eval": "Overweight",
                        "grading": 5,
                        "message": "Being overweight can increase your risk of heart disease and diabetes",
                        "refer": "no",
                        "target": "18.5-23",
                        "tfl": "AMBER",
                        "value": "24.71"
                    },
                    "body_fat": {
                        "code": "BFT-OVERWEIGHT",
                        "eval": "Overweight",
                        "grading": 5,
                        "message": "Body fat is in the obese range. Discuss with health professional",
                        "refer": "no",
                        "target": "13%-25%",
                        "tfl": "AMBER",
                        "value": "25.0%"
                    },
                    "waist_circ": {
                        "code": "WST-NORMAL",
                        "eval": "Healthy level",
                        "grading": 0,
                        "message": "Waist circumference is within healthy range",
                        "refer": "no",
                        "target": "101cm",
                        "tfl": "GREEN",
                        "value": "100.7cm"
                    },
                    "whr": {
                        "code": "WHR-HIGH",
                        "eval": "Abnormal",
                        "grading": 10,
                        "message": "Your waist hip ratio is above the target",
                        "refer": "no",
                        "target": "0.90",
                        "tfl": "AMBER",
                        "value": "0.92"
                    }
                },
                "message": "You are in abnormal range. Try to do more exercise and follow guidelines for healthy life."
            },
            "cholesterol": {
                "components": {
                    "hdl": {
                        "code": "",
                        "eval": "",
                        "grading": 0,
                        "message": "",
                        "refer": "",
                        "target": "",
                        "tfl": "",
                        "value": ""
                    },
                    "ldl": {
                        "code": "",
                        "eval": "",
                        "grading": 0,
                        "message": "",
                        "refer": "",
                        "target": "",
                        "tfl": "",
                        "value": ""
                    },
                    "tg": {
                        "code": "",
                        "eval": "",
                        "grading": 0,
                        "message": "",
                        "refer": "",
                        "target": "",
                        "tfl": "",
                        "value": ""
                    },
                    "total_cholesterol": {
                        "code": "CHOL-ELEVATED-WITH-HIGH-CVD-RISK",
                        "eval": "Elevated",
                        "grading": 0,
                        "message": "Based on Tchol > 5 and high CVD risk, requires medication",
                        "refer": "yes",
                        "target": "< 90mg/dL (5 mmol/L)",
                        "tfl": "AMBER",
                        "value": "5.6mmol/L"
                    }
                },
                "message": ""
            },
            "cvd": {
                "code": "CVD-HIGH-RISK",
                "eval": "",
                "grading": 0,
                "message": "",
                "refer": "",
                "target": "Risk Management",
                "tfl": "",
                "value": "20-30%"
            },
            "diabetes": {
                "code": "DM-EXISTING-POOR-CONTROL",
                "eval": "Poor Control",
                "grading": 0,
                "message": "Diabetes is poorly controlled. Review",
                "refer": "yes",
                "target": "<117mg/dL (6.5mmol/L)",
                "tfl": "AMBER",
                "value": "15.6mmol/L"
            },
            "lifestyle": {
                "components": {
                    "alcohol": {
                        "code": "ALC-LOW-RISK",
                        "eval": "Low Risk",
                        "grading": 0,
                        "message": "If you regularly drink 14 units per week, it's best to spread your drinking over 3 or more days (avoid heavy drinking on any single day)",
                        "refer": "no",
                        "target": "< 14 units/week",
                        "tfl": "AMBER",
                        "value": "12.0 units"
                    },
                    "diet": {
                        "components": {
                            "fruit": {
                                "code": "FRT-TARGET",
                                "eval": "On target",
                                "grading": 0,
                                "message": "Keep eating fruit!",
                                "refer": "no",
                                "target": ">2 servings/day",
                                "tfl": "GREEN",
                                "value": "21 servings"
                            },
                            "vegetable": {
                                "code": "VEG-LOW",
                                "eval": "Below target",
                                "grading": 5,
                                "message": "Eat more vegetables",
                                "refer": "no",
                                "target": "5 or more servings/day",
                                "tfl": "AMBER",
                                "value": "14 servings"
                            }
                        },
                        "message": "You are not following healthy diet. Please make sure to consume ample fruits and vegetables."
                    },
                    "physical_activity": {
                        "code": "PA-UNDER",
                        "eval": "Below target",
                        "grading": 0,
                        "message": "Not quit meeting targets. Aim for 150 minutes of moderate OR 75 minutes of vigroous intensity exercise per week",
                        "refer": "no",
                        "target": ">150 minutes",
                        "tfl": "AMBER",
                        "value": "140 minutes"
                    },
                    "smoking": {
                        "code": "SM-EX-SMOKER",
                        "eval": "Non Smoker",
                        "grading": 0,
                        "message": "Good work quitting. It's the best thing you can do for your health. Stay strong",
                        "refer": "no",
                        "target": "no",
                        "tfl": "GREEN",
                        "value": "ex smoker"
                    }
                },
                "message": ""
            }
        },
        "goals": [
            {
                "code": "BSC-1",
                "eval": "Diabetes Management",
                "tfl": "AMBER",
                "message": "Consult and take appripriate actions by reviewing with a doctor to get your blood sugar level normal."
            },
            {
                "code": "WC-1",
                "eval": "Weight control",
                "tfl": "RED",
                "message": "Need to reduce weight. Consult with a dietician."
            },
            {
                "code": "MC-2",
                "eval": "CVD Risk Review",
                "tfl": "AMBER",
                "message": "You should review with doctor as you have high CVD risk."
            }
        ],
        "meta": {
            "algorithm": "Hearts Algorithm",
            "request_id": "fcdff736-5178-40d1-8c3a-727b057400bd"
        },
        "referrals": {
            "reasons": [
                {
                    "type": "diabetes",
                    "urgent": false
                },
                {
                    "type": "blood pressure",
                    "urgent": false
                },
                {
                    "type": "cvd",
                    "urgent": false
                },
                {
                    "type": "total cholesterol",
                    "urgent": false
                }
            ],
            "refer": true,
            "urgent": false
        }
    }
}
```

## Algorithm Input Parameters

The inputs can be divided into many parts. Let us discuss those.

### config

First part is **config**. It consist of the algorithm name and the risk model to be used with the algorithm.

```javascript
"config": {
    "algorithm": "hearts",
    "risk_model": "whocvd"
}
```

### param

All parameters are inside this object.

#### demographics

It consists of patient/participant's demographical information.
Mandatory elements: **gender**, **age**, **birth_country_code**

```javascript
"demographics": {
    "gender": "male",
    "age": {
        "value": 65,
        "unit": "year"
    },
    "birth_country": "Bangladesh",
    "birth_country_code": "BD",
    "living_country": "Bangladesh",
    "living_country_code": "BD",
    "race": "Bengali",
    "ethnicity": "Bengali"
}
```

#### components/lifestyle

All the lifestyle components are grouped inside an array.

##### smoking

Possible **value**: smoker, non-smoker, ex-smoker.
If ex-smoker and quitted within a year, set quit_within_year to **true**.

```javascript
{
    "name": "smoking",
    "category": "addiction",
    "value": "ex-smoker",
    "quit_within_year": true
}
```

##### alcohol

Possible **frequency**: daily, weekly, monthly.
To calculate the units of consumption, use the following method to calculate:

[Calculating Units >](https://www.nhs.uk/live-well/alcohol-support/calculating-alcohol-units/#calculating-units)

```javascript
{
    "name": "alcohol_history",
    "category": "addiction",
    "value": 12,
    "units": "units",
    "frequency": "weekly"
}
```

##### fruit

Possible **frequency**: daily, weekly.

```javascript
{
    "name": "fruit",
    "category": "nutrition",
    "value": 3,
    "units": "servings",
    "frequency": "daily"
}
```

##### vegetables

Possible **frequency**: daily, weekly.

```javascript
{
    "name": "vegetables",
    "category": "nutrition",
    "value": 2,
    "units": "servings",
    "frequency": "daily"
}
```

##### exercise

Possible **frequency**: daily, weekly.
Possible **name**: low-intensity-exercise, high-intensity-exercise.

```javascript
{
    "name": "low-intensity-exercise",
    "category": "physical-activity",
    "value": 20,
    "units": "minutes",
    "frequency": "daily"
}
```

#### components/body-measurements

Body measurements input params are added in this array.

##### blood_pressure

Possible **arm**: right, left.
**MULTIPLE** inputs.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "blood_pressure",
    "category": "vital-signs",
    "value": "170/80",
    "units": "mmHg",
    "arm": "right"
}
```

##### height/hip/waist

Possible **name**: height, hip, waist.
Possible **units**: cm, m, inch.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "height",
    "category": "anthropometry",
    "units": "cm",
    "value": 180.5
}
```

##### weight

Possible **units**: kg, lb.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "weight",
    "category": "anthropometry",
    "units": "kg",
    "value": 80.5
}
```

##### body-fat/muscle/visceral_fat

Possible **name**: body-fat, muscle, visceral_fat.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "body-fat",
    "category": "anthropometry",
    "units": "%",
    "value": 25
}
```

#### components/biological-samples

Biological samples includes blood samples for blood sugar and cholesterol.

##### blood_sugar/a1c

Only one of this input should be in the params.

The **a1c** input will have the following structure.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "a1c",
    "category": "blood-test",
    "value": 6.5,
    "units": "%"
}
```

The **blood_sugar** input have little bit different structure.
Possible **units**: mmol/L, mg/dL. Possible **type**: fasting, random.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "blood_sugar",
    "category": "blood-test",
    "value": 15.6,
    "units": "mmol/L",
    "type": "fasting"
}
```

##### cholesterol

Possible **name**: total_cholesterol, hdl, ldl, tg.
Possible **units**: mmol/L, mg/dL.

```javascript
{
    "effectiveDate": "timestamp",
    "name": "total_cholesterol",
    "category": "blood-test",
    "value": 5.6,
    "units": "mmol/L"
}
```

#### components/medical_history

Add your conditions as objects inside the array.

Possible **is_active**: true, false.

```javascript
{
    "name": "hypertension",
    "category": "condition",
    "is_active": true
}
```

You should add the allergies in this array too.

```javascript
{
    "category": "allergy",
    "type": "medication",
    "allergen": "ace-i",
    "criticality": "",
    "reaction": "cough"
}
```

#### components/medications

List all the current medications in this array.

```javascript
{
    "generic": "rampiril",
    "category": "anti-hypertensive"
}
```

#### components/family_history

List all the high level diseases of the relatives.

```javascript
{
    "name": "cardiovascular-disease",
    "relative": "1st degree"
}
```

## Algorithm Output

The algorithm output is also divided into many sections.

### hearts

All the hearts algorithm output is grouped inside this object.

#### assessments

All assessments are inside this object. They are either direct assessment of any attribute or grouped inside `components`.

Possible **tfl**: DARK RED, RED, AMBER, YELLOW, ORANGE, GREEN.

```javascript
{
    "code": "BP-HTN",
    "eval": "Hypertension",
    "grading": 0,
    "message": "Your BP indicates hypertension. Discuss with local health professional",
    "refer": "yes",
    "target": "130/80",
    "tfl": "AMBER",
    "value": "145/85"
}
```

We have the following list of assessments:

- blood_pressure
- body_composition/components/bmi
- body_composition/components/body_fat
- body_composition/components/waist_circ
- body_composition/components/whr
- cholesterol/components/ldl
- cholesterol/components/hdl
- cholesterol/components/tg
- cholesterol/components/total_cholesterol
- cvd
- diabetes
- lifestyle/components/alcohol
- lifestyle/components/diet/components/fruit
- lifestyle/components/diet/components/vegetable
- lifestyle/components/physical_activity
- lifestyle/components/smoking

#### goals

It will consists of all the generated goals.

```javascript
{
    "code": "BSC-1",
    "eval": "Diabetes Management",
    "tfl": "AMBER",
    "message": "Consult and take appripriate actions by reviewing with a doctor to get your blood sugar level normal."
}
```

#### meta

Meta information for each request will be in this section.

```javascript
"meta": {
    "algorithm": "Hearts Algorithm",
    "request_id": "c94fea6a-91ea-4f17-97c3-2455aeb3fdb1"
}
```

#### referrals

It will show all the referral reasons as well as overall referral and urgency.

```javascript
"referrals": {
    "reasons": [
        {
            "type": "diabetes",
            "urgent": false
        },
        {
            "type": "blood pressure",
            "urgent": false
        },
        {
            "type": "cvd",
            "urgent": false
        },
        {
            "type": "total cholesterol",
            "urgent": false
        }
    ],
    "refer": true,
    "urgent": false
}
```

## Errors

The response will show any error related to assesments. In this case, the response code is 200 but there will be error in the `errors` array.

``` javascript
{
    "errors": [
        "no matching condition found for whr"
    ],
    "hearts: {
        ...
    }
}
```

For any invalid request, it will return the error code (like 400, 403, 404, 405, 422 etc..) along with a error object.

```javascript
{
    "error": "invalid method, only accepts post request"
}
```

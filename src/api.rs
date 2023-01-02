use std::{fmt::Display, net::TcpStream};

use thiserror::Error;

use url::Url;

use crate::io_models::{
    JeedomBatteryResult, JeedomGlobalSummaryResult, JeedomJsonRpcRequest, JeedomNotificationResult,
    JeedomPingResult, JeedomRpcError, JeedomVersionResult,
};

#[derive(Debug, PartialEq)]
pub struct JeedomClient {
    pub url: String,
    apikey: String,
    debug: bool,
}

#[derive(Error, Debug, PartialEq)]
pub enum ConnectivityErrors {
    #[error("URL can't be parsed: {0}")]
    UrlNotParsable(String),
    #[error("URL not reachable")]
    NoGivenUrlCanBeReached,
    #[error("unable to get the host from the given URL: {0}")]
    UrlHostUnknown(String),
    #[error("unable to get the port from the given URL: {0}")]
    UrlPortUnknown(String),
}

pub enum JeedomClientErrors {
    ApiError(Box<ureq::Error>),
    DeserializationError(Box<serde_json::Error>),
}

impl Display for JeedomClientErrors {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::ApiError(e) => {
                match serde_json::from_str::<JeedomRpcError>(e.to_string().as_str()) {
                    Ok(e) => write!(f, "Jeedom API error: {}", e.error.message),
                    Err(e) => write!(f, "unable to parse the error: {e}"),
                }
            }
            Self::DeserializationError(e) => {
                match serde_json::from_str::<JeedomRpcError>(
                    e.to_string().replace('\"', "").as_str(),
                ) {
                    Ok(e) => write!(f, "Deserialization error: {}", e.error.message),
                    Err(e) => write!(f, "unable to parse the error: {e}"),
                }
            }
        }
    }
}

impl JeedomClient {
    pub fn new(jeedom_url: String, apikey: String, debug: bool) -> Self {
        Self {
            url: format!("{jeedom_url}/core/api/jeeApi.php"),
            apikey,
            debug,
        }
    }

    pub fn check_connectivity(
        local_url: String,
        external_url: Option<String>,
        debug: bool,
    ) -> Result<String, ConnectivityErrors> {
        let mut urls = vec![local_url];
        if let Some(external_url_present) = external_url {
            urls.push(external_url_present);
        }

        for url in urls {
            if debug {
                print!("Checking {url} connectivity...");
            }
            let url_format = Url::parse(url.as_str())
                .map_err(|e| ConnectivityErrors::UrlNotParsable(e.to_string()))?;
            let host = match url_format.host_str() {
                Some(x) => x,
                None => return Err(ConnectivityErrors::UrlHostUnknown(url)),
            };
            let port = match url_format.port_or_known_default() {
                Some(x) => x,
                None => return Err(ConnectivityErrors::UrlPortUnknown(url)),
            };
            if TcpStream::connect(format!("{host}:{port}").as_str()).is_ok() {
                if debug {
                    println!("OK");
                }
                return Ok(url);
            }
        }

        Err(ConnectivityErrors::NoGivenUrlCanBeReached)
    }

    fn api_request(&self, action: String) -> Result<String, Box<ureq::Error>> {
        let params = JeedomJsonRpcRequest::new(self.apikey.clone(), action.clone());
        let request_result = ureq::post(self.url.as_str()).send_json(ureq::json!(&params))?;
        let result_string = request_result
            .into_string()
            .unwrap_or_else(|_| panic!("couldn't convert {action} request result to string"))
            .replace('\t', "");
        if self.debug {
            // no need tracing for this simple debug message
            println!("\nDEBUG: API request '{action}' result:\n{result_string:?}");
        }

        Ok(result_string)
    }

    pub fn ping(&self) -> Result<JeedomPingResult, JeedomClientErrors> {
        let body: JeedomPingResult = serde_json::from_str(
            self.api_request("ping".to_string())
                .map_err(JeedomClientErrors::ApiError)?
                .as_str(),
        )
        .map_err(|e| JeedomClientErrors::DeserializationError(Box::new(e)))?;

        Ok(body)
    }

    pub fn global_summary(&self) -> Result<JeedomGlobalSummaryResult, JeedomClientErrors> {
        let body: JeedomGlobalSummaryResult = serde_json::from_str(
            self.api_request("summary::global".to_string())
                .map_err(JeedomClientErrors::ApiError)?
                .as_str(),
        )
        .map_err(|e| JeedomClientErrors::DeserializationError(Box::new(e)))?;

        Ok(body)
    }

    pub fn battery_status(&self) -> Result<JeedomBatteryResult, JeedomClientErrors> {
        let body: JeedomBatteryResult = serde_json::from_str(
            self.api_request("eqLogic::all".to_string())
                .map_err(JeedomClientErrors::ApiError)?
                .as_str(),
        )
        .map_err(|e| JeedomClientErrors::DeserializationError(Box::new(e)))?;

        Ok(body)
    }

    pub fn notification_messages(&self) -> Result<JeedomNotificationResult, JeedomClientErrors> {
        let body: JeedomNotificationResult = serde_json::from_str(
            self.api_request("message::all".to_string())
                .map_err(JeedomClientErrors::ApiError)?
                .as_str(),
        )
        .map_err(|e| JeedomClientErrors::DeserializationError(Box::new(e)))?;

        Ok(body)
    }

    pub fn version(&self) -> Result<JeedomVersionResult, JeedomClientErrors> {
        let body: JeedomVersionResult = serde_json::from_str(
            self.api_request("version".to_string())
                .map_err(JeedomClientErrors::ApiError)?
                .as_str(),
        )
        .map_err(|e| JeedomClientErrors::DeserializationError(Box::new(e)))?;

        Ok(body)
    }
}

#[cfg(test)]
mod tests {
    use std::{fs, net::TcpListener};

    use crate::{
        api::ConnectivityErrors,
        io_models::{
            JeedomBatteryResult, JeedomGlobalSummaryResult, JeedomNotificationResult,
            JeedomRpcError,
        },
    };

    use super::JeedomClient;

    #[test]
    fn jeedom_client_new() {
        let url = "http://127.0.0.1:8080".to_string();
        let apikey = "xxxxxx".to_string();
        assert_eq!(
            JeedomClient::new(url.clone(), apikey.clone(), false),
            JeedomClient {
                url: format!("{url}/core/api/jeeApi.php"),
                apikey,
                debug: false
            }
        )
    }

    #[test]
    fn test_check_connectivity() {
        let local_addr = "127.0.0.1:8080";
        let local_url = format!("http://{local_addr}");
        let external_addr = "127.0.0.1:8081";
        let external_url = format!("http://{external_addr}");
        let bad_host = ":8082";
        let bad_host_url = format!("http://{bad_host}");
        let bad_port = "xxx:80822";
        let bad_port_url = format!("http://{bad_port}");

        // check bad url format
        assert_eq!(
            JeedomClient::check_connectivity(bad_host_url, None, false),
            Err(ConnectivityErrors::UrlNotParsable("empty host".to_string()))
        );
        assert_eq!(
            JeedomClient::check_connectivity(bad_port_url, None, false),
            Err(ConnectivityErrors::UrlNotParsable(
                "invalid port number".to_string()
            ))
        );

        // check when nothing works
        assert_eq!(
            JeedomClient::check_connectivity(local_url.clone(), None, false),
            Err(ConnectivityErrors::NoGivenUrlCanBeReached)
        );
        assert_eq!(
            JeedomClient::check_connectivity(local_url.clone(), Some(external_url.clone()), false),
            Err(ConnectivityErrors::NoGivenUrlCanBeReached)
        );

        // only external url is working
        let _external_listener = TcpListener::bind(external_addr).unwrap();
        assert_eq!(
            JeedomClient::check_connectivity(local_url.clone(), Some(external_url.clone()), false)
                .unwrap(),
            external_url
        );

        // first url is working
        let _local_listener = TcpListener::bind("127.0.0.1:8080").unwrap();
        assert_eq!(
            JeedomClient::check_connectivity(local_url.clone(), Some(external_url), false).unwrap(),
            local_url
        );
    }

    #[test]
    fn global_status_deserialize() {
        let request_result = r#"
{
    "jsonrpc": "2.0",
    "id": "1",
    "result": {
        "security": {
            "key": "security",
            "name": "Alerte",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-alerte2\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 0
        },
        "motion": {
            "key": "motion",
            "name": "Mouvement",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-mouvement\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 0
        },
        "door": {
            "key": "door",
            "name": "Porte",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-porte-ouverte\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 1
        },
        "windows": {
            "key": "windows",
            "name": "Fenu00eatre",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-fenetre-ouverte\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 0
        },
        "shutter": {
            "key": "shutter",
            "name": "Volet",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-volet-ouvert\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": null
        },
        "light": {
            "key": "light",
            "name": "Lumiu00e8re",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-lumiere-on\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 0
        },
        "outlet": {
            "key": "outlet",
            "name": "Prise",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-prise\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 2
        },
        "temperature": {
            "key": "temperature",
            "name": "Tempu00e9rature",
            "calcul": "avg",
            "icon": "<i class=\"icon divers-thermometer31\"></i>",
            "iconnul": "",
            "unit": "u00b0C",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "",
            "allowDisplayZero": "1",
            "ignoreIfCmdOlderThan": "",
            "value": null
        },
        "humidity": {
            "key": "humidity",
            "name": "Humiditu00e9",
            "calcul": "avg",
            "icon": "<i class=\"fa fa-tint\"></i>",
            "iconnul": "",
            "unit": "%",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "",
            "allowDisplayZero": "1",
            "ignoreIfCmdOlderThan": "",
            "value": null
        },
        "luminosity": {
            "key": "luminosity",
            "name": "Luminositu00e9",
            "calcul": "avg",
            "icon": "<i class=\"icon meteo-soleil\"></i>",
            "iconnul": "",
            "unit": "lx",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": null
        },
        "power": {
            "key": "power",
            "name": "Puissance",
            "calcul": "sum",
            "icon": "<i class=\"fa fa-bolt\"></i>",
            "iconnul": "",
            "unit": "W",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": null
        },
        "alarm": {
            "key": "alarm",
            "name": "Alarme",
            "calcul": "sum",
            "icon": "<i class=\"icon jeedom-lock-ferme\"></i>",
            "iconnul": "",
            "unit": "",
            "hidenumber": "0",
            "hidenulnumber": "0",
            "count": "binary",
            "allowDisplayZero": "0",
            "ignoreIfCmdOlderThan": "",
            "value": 0
        }
    }
}
    "#;
        let x = serde_json::from_str(request_result);
        assert!(x.is_ok());
        let _: JeedomGlobalSummaryResult = x.unwrap();
    }

    #[test]
    #[ignore]
    fn battery_status_deserialize_from_file() {
        let request_result = fs::read_to_string("json.json").unwrap();
        let x = serde_json::from_str::<JeedomBatteryResult>(request_result.as_str());
        println!("{x:?}");
        assert!(x.is_ok());
    }

    #[test]
    fn battery_status_deserialize() {
        let request_result = r#"
{
    "id": "1",
    "jsonrpc": "2.0",
    "result": [
        {
            "cache": {
                "waiting": []
            },
            "category": {
                "automatism": "0",
                "default": "0",
                "energy": "0",
                "heating": "0",
                "light": "0",
                "multimedia": "0",
                "opening": "0",
                "security": "0"
            },
            "comment": null,
            "configuration": {
                "createtime": "2023-01-30 10:39:21",
                "firmwareVersion": "1.0",
                "interview": "complete",
                "manufacturer_id": 134,
                "product_id": 90,
                "product_name": "ZW090 - Zu2010Stick Gen5 USB Controller",
                "product_type": 1,
                "refreshes": [],
                "updatetime": "2023-01-30 11:53:12"
            },
            "display": {
                "backGraph::info": 0
            },
            "eqType_name": "zwavejs",
            "generic_type": null,
            "id": "576",
            "isEnable": "1",
            "isVisible": "0",
            "logicalId": "1",
            "name": "1 - AEON Labs Zu2010Stick Gen5 USB Controller ZW090",
            "object_id": null,
            "order": "9999",
            "status": {
                "danger": 0,
                "lastCommunication": "2023-01-31 08:43:32",
                "timeout": 0,
                "warning": 0
            },
            "tags": null,
            "timeout": null
        },
        {
            "cache": [],
            "category": {
                "automatism": "0",
                "default": "1",
                "energy": "0",
                "heating": "0",
                "light": "0",
                "multimedia": "1",
                "opening": "0",
                "security": "0"
            },
            "comment": null,
            "configuration": {
                "IPV4": "0",
                "IPV6": "0",
                "VersionLogicalID": "2.1",
                "action": "",
                "autorefresh": "*/5 * * * *",
                "createtime": "2021-09-15 09:10:14",
                "eq_group": "system",
                "info": "",
                "logicalID": "airmedia",
                "type": "",
                "type2": "",
                "updatetime": "2022-03-17 09:55:19"
            },
            "display": [],
            "eqType_name": "Freebox_OS",
            "generic_type": null,
            "id": "382",
            "isEnable": "0",
            "isVisible": "0",
            "logicalId": "airmedia",
            "name": "Air Mu00e9dia",
            "object_id": null,
            "order": "9999",
            "status": [],
            "tags": null,
            "timeout": null
        },
        {
            "cache": {
                "waiting": []
            },
            "category": {
                "automatism": "0",
                "default": "0",
                "energy": "0",
                "heating": "1",
                "light": "0",
                "multimedia": "0",
                "opening": "0",
                "security": "0"
            },
            "comment": null,
            "configuration": {
                "battery_type": "2x1.5V AA",
                "createtime": "2023-01-30 10:39:22",
                "firmwareVersion": "0.16",
                "interview": "complete",
                "manufacturer_id": 328,
                "product_id": 1,
                "product_name": "Spirit - Thermostatic Valve",
                "product_type": 3,
                "refreshes": [],
                "updatetime": "2023-01-30 23:54:42"
            },
            "display": {
                "backGraph::info": 0,
                "height": "474px",
                "width": "232px"
            },
            "eqType_name": "zwavejs",
            "generic_type": null,
            "id": "578",
            "isEnable": "1",
            "isVisible": "1",
            "logicalId": "22",
            "name": "22 - Radiateur",
            "object_id": "12",
            "order": "11",
            "status": {
                "battery": 70,
                "batteryDatetime": "2023-01-31 08:44:49",
                "batterydanger": 0,
                "batterywarning": 0,
                "danger": 0,
                "lastCommunication": "2023-01-31 13:27:06",
                "timeout": 0,
                "warning": 0
            },
            "tags": null,
            "timeout": null
        },
        {
            "cache": {
                "waiting": []
            },
            "category": {
                "automatism": "0",
                "default": "0",
                "energy": "0",
                "heating": "0",
                "light": "0",
                "multimedia": "0",
                "opening": "1",
                "security": "1"
            },
            "comment": null,
            "configuration": {
                "battery_type": "1x3.6V ER14250",
                "batterytime": "2023-01-30 13:36:30",
                "createtime": "2023-01-30 10:47:21",
                "firmwareVersion": "3.2",
                "interview": "complete",
                "lastWakeUp": 1675158665,
                "manufacturer_id": 271,
                "product_id": 4096,
                "product_name": "FGDW002 - Fibaro Door Window Sensor 2",
                "product_type": 1794,
                "refreshes": [],
                "updatetime": "2023-01-31 10:51:05"
            },
            "display": {
                "backGraph::info": 0,
                "height": "200px",
                "width": "232px"
            },
            "eqType_name": "zwavejs",
            "generic_type": null,
            "id": "618",
            "isEnable": "1",
            "isVisible": "1",
            "logicalId": "39",
            "name": "39 - Fenetre",
            "object_id": "12",
            "order": "10",
            "status": {
                "battery": 100,
                "batteryDatetime": "2023-01-31 08:43:33",
                "batterydanger": 0,
                "batterywarning": 0,
                "danger": 0,
                "lastCommunication": "2023-01-31 13:05:53",
                "timeout": 0,
                "warning": 0
            },
            "tags": null,
            "timeout": null
        }
    ]
}
        "#;
        let x = serde_json::from_str(request_result);
        assert!(x.is_ok());
        let _: JeedomBatteryResult = x.unwrap();
    }

    #[test]
    fn notification_deserialize() {
        let request_result = r#"
        {
            "id": "1",
            "jsonrpc": "2.0",
            "result": [
                {
                    "action": "",
                    "date": "2023-02-10 16:38:02",
                    "id": "2666",
                    "logicalId": "update::xxx",
                    "message": "De nouvelles mises u00e0 jour sont disponibles : Freebox_OS",
                    "occurrences": null,
                    "plugin": "update"
                },
                {
                    "action": "",
                    "date": "2023-02-09 07:17:41",
                    "id": "2665",
                    "logicalId": "mobile::yyy",
                    "message": "Erreur exu00e9cution de la commande [Parcelle][Pixel 5][Notification] : Echec de l&apos;envoi de la notification :{&quot;state&quot;:&quot;nok&quot;,&quot;error&quot;:&quot;{&quot;code&quot;:&quot;messaging/invalid-argument&quot;,&quot;message&quot;:&quot;The registration token is not a valid FCM registration token&quot;}&quot;}",
                    "occurrences": "3",
                    "plugin": "mobile"
                }
            ]
        }
        "#;

        let x = serde_json::from_str(request_result);
        assert!(x.is_ok());
        let _: JeedomNotificationResult = x.unwrap();
    }

    #[test]
    fn test_deserialize_error() {
        let request_result = r#"
        {
            "jsonrpc":"2.0",
            "id":99999,
            "error":
                {
                    "code":2,
                    "message":"Vous n'u00eates pas autorisu00e9 u00e0 effectuer cette action"
                }
        }
        "#;
        let x = serde_json::from_str(request_result);
        assert!(x.is_ok());
        let err: JeedomRpcError = x.unwrap();
        assert_eq!(err.error.code, 2);
        println!("err: {err:?}");
    }
}

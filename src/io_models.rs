use std::{
    fmt::{self},
    marker::PhantomData,
};

use derivative::Derivative;
use serde::{
    de::{self, MapAccess, SeqAccess, Visitor},
    Deserialize, Deserializer, Serialize,
};
use serde_with::serde_as;

//
// Jeedom Json RPC input and errors
//

#[derive(Serialize, Deserialize)]
pub struct JeedomJsonRpcRequest {
    pub jsonrpc: String,
    pub id: String,
    pub method: String,
    pub params: JeedomJsonRpcApiParamsRequest,
}

#[derive(Serialize, Deserialize, Derivative)]
#[derivative(Debug)]
pub struct JeedomJsonRpcApiParamsRequest {
    #[derivative(Debug = "ignore")]
    apikey: String,
    datetime: String,
}

impl JeedomJsonRpcRequest {
    pub fn new(apikey: String, method: String) -> Self {
        Self {
            jsonrpc: "2.0".to_string(),
            id: "1".to_string(),
            method,
            params: JeedomJsonRpcApiParamsRequest {
                apikey,
                datetime: "1".to_string(),
            },
        }
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomRpcError {
    pub error: JeedomErrorContent,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomErrorContent {
    pub code: u32,
    pub message: String,
}

//
// Jeedom results
//

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomPingResult {
    pub result: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomVersionResult {
    pub result: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomGlobalSummaryResult {
    pub result: JeedomGlobalSummary,
}
#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomGlobalSummary {
    pub security: JeedomSummary,
    pub motion: JeedomSummary,
    pub door: JeedomSummary,
    pub windows: JeedomSummary,
    pub shutter: JeedomSummary,
    pub light: JeedomSummary,
    pub outlet: JeedomSummary,
    pub temperature: JeedomSummary,
    pub humidity: JeedomSummary,
    pub luminosity: JeedomSummary,
    pub power: JeedomSummary,
    pub alarm: Option<JeedomSummary>, // my custom summary
}

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomSummary {
    pub key: String,
    pub name: String,
    pub unit: String,
    pub value: Option<u32>,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomBatteryResult {
    pub result: Vec<JeedomBattery>,
}
#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomBattery {
    #[serde(deserialize_with = "deserialize_vec_or_struct")]
    pub status: Option<JeedomBatteryStatus>,
}

#[serde_as]
#[derive(Serialize, Deserialize, Debug)]
pub struct JeedomBatteryStatus {
    pub batterydanger: Option<u32>,
    pub batterywarning: Option<u32>,
    pub battery: Option<BatteryPercentage>,
}

#[derive(Serialize, Debug)]
pub struct BatteryPercentage(pub Option<u32>);

impl<'de> Deserialize<'de> for BatteryPercentage {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
    where
        D: Deserializer<'de>,
    {
        struct MyVisitor;

        impl<'de> Visitor<'de> for MyVisitor {
            type Value = BatteryPercentage;

            fn expecting(&self, fmt: &mut fmt::Formatter<'_>) -> fmt::Result {
                fmt.write_str("integer or string")
            }

            fn visit_u64<E>(self, val: u64) -> Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                let x = val as u32;
                Ok(BatteryPercentage(Some(x)))
            }

            fn visit_str<E>(self, val: &str) -> Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match val.parse::<u32>() {
                    Ok(val) => self.visit_u32(val),
                    Err(_) => Ok(BatteryPercentage(None)),
                }
            }
        }

        deserializer.deserialize_any(MyVisitor)
    }
}

// need this visitor to support status[] and status{} from JeedomBatteryStatus
fn deserialize_vec_or_struct<'de, T, D>(deserializer: D) -> Result<Option<T>, D::Error>
where
    T: Deserialize<'de>,
    D: Deserializer<'de>,
{
    struct DeserializeVecOrStructBatteryStatus<T>(PhantomData<fn() -> Option<T>>);

    impl<'de, T> de::Visitor<'de> for DeserializeVecOrStructBatteryStatus<T>
    where
        T: Deserialize<'de>,
    {
        type Value = Option<T>;

        fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
            formatter.write_str("a vec or a struct")
        }

        fn visit_seq<V>(self, _: V) -> Result<Self::Value, V::Error>
        where
            V: SeqAccess<'de>,
        {
            Ok(None)
        }

        fn visit_map<S>(self, jeedom_battery: S) -> Result<Self::Value, S::Error>
        where
            S: MapAccess<'de>,
        {
            let x = Deserialize::deserialize(de::value::MapAccessDeserializer::new(jeedom_battery))
                .map_err(de::Error::custom)?;
            Ok(Some(x))
        }
    }

    deserializer.deserialize_any(DeserializeVecOrStructBatteryStatus(PhantomData))
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct JeedomUpdateStatusResult {
    pub result: Vec<JeedomUpdateStatus>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct JeedomUpdateStatus {
    pub name: String,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct JeedomNotificationResult {
    pub result: Vec<JeedomNotificationStatus>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct JeedomNotificationStatus {
    pub plugin: String,
}

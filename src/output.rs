use crate::{
    io_models::{
        BatteryPercentage, JeedomBatteryResult, JeedomGlobalSummaryResult, JeedomNotificationResult,
    },
    BarStyle, BarType,
};
use colored::Colorize;
use phf::phf_map;

// hard limit danger if nothing is set
const JEEDOM_BATTERY_MIN_DANGER: u32 = 5;
const JEEDOM_BATTERY_MIN_WARNING: u32 = 20;
// jeedom notification icons
static JEEDOM_NOTIFICATION_ICONS: phf::Map<&'static str, &'static str> = phf_map! {
    "1" => "\u{2460}",
    "2" => "\u{2461}",
    "3" => "\u{2462}",
    "4" => "\u{2463}",
    "5" => "\u{2464}",
    "6" => "\u{2465}",
    "7" => "\u{2466}",
    "8" => "\u{2467}",
    "9" => "\u{2468}",
    "10" => "\u{2469}",
    "11" => "\u{246A}",
    "12" => "\u{246B}",
    "13" => "\u{246C}",
    "14" => "\u{246D}",
    "15" => "\u{246E}",
    "16" => "\u{246F}",
    "17" => "\u{2470}",
    "18" => "\u{2471}",
    "19" => "\u{2472}",
    "20" => "\u{2473}",
};

// icons
pub struct JeedomSummaryOutput {
    pub alarm: JeedomSummaryItemOutput,
    pub battery: JeedomBatteryInfoOuput,
    pub door: JeedomSummaryItemOutput,
    pub humidity: JeedomSummaryItemOutput,
    pub light: JeedomSummaryItemOutput,
    pub luminosity: JeedomSummaryItemOutput,
    pub motion: JeedomSummaryItemOutput,
    pub outlet: JeedomSummaryItemOutput,
    pub power: JeedomSummaryItemOutput,
    pub security: JeedomSummaryItemOutput,
    pub shutter: JeedomSummaryItemOutput,
    pub temperature: JeedomSummaryItemOutput,
    pub updates: JeedomSummaryItemOutput,
    pub windows: JeedomSummaryItemOutput,
    pub notifications: JeedomNotificationsItemOutput,
}

pub struct JeedomSummaryItemOutput {
    pub icon: String,
    pub counter: u32,
}

impl JeedomSummaryItemOutput {
    pub fn new(icon: String) -> Self {
        Self { icon, counter: 0 }
    }
}
#[derive(Default, Clone)]
pub struct JeedomBatteryInfoOuput {
    pub icon: String,
    pub battery_warning_counter: u32,
    pub battery_danger_counter: u32,
}

impl JeedomBatteryInfoOuput {
    pub fn new(icon: String) -> Self {
        Self {
            icon,
            battery_warning_counter: 0,
            battery_danger_counter: 0,
        }
    }
}
#[derive(Default)]
pub struct JeedomNotificationsItemOutput {
    pub message_icon: Option<String>,
    pub message_counter: u32,
    pub update_icon: Option<String>,
    pub update_counter: u32,
}

impl JeedomSummaryOutput {
    pub fn new(
        global_summary: JeedomGlobalSummaryResult,
        battery_status: JeedomBatteryResult,
        notifications: JeedomNotificationResult,
        bar_style: BarStyle,
    ) -> Self {
        // setup icons
        let mut jeedom_summary_output = match bar_style {
            BarStyle::Nerd => JeedomSummaryOutput::init_with_nerd_fonts_icons(),
            BarStyle::Emoji => JeedomSummaryOutput::init_with_emoji_icons(),
            BarStyle::Jeedom => JeedomSummaryOutput::init_with_jeedom_icons(),
            BarStyle::Text => JeedomSummaryOutput::init_without_icons(),
        };

        // update global summary counters
        if let Some(x) = global_summary.result.alarm {
            jeedom_summary_output.alarm.counter = x.value.unwrap_or(0);
        };
        jeedom_summary_output.door.counter = global_summary.result.door.value.unwrap_or(0);
        jeedom_summary_output.humidity.counter = global_summary.result.humidity.value.unwrap_or(0);
        jeedom_summary_output.light.counter = global_summary.result.light.value.unwrap_or(0);
        jeedom_summary_output.luminosity.counter =
            global_summary.result.luminosity.value.unwrap_or(0);
        jeedom_summary_output.motion.counter = global_summary.result.motion.value.unwrap_or(0);
        jeedom_summary_output.outlet.counter = global_summary.result.outlet.value.unwrap_or(0);
        jeedom_summary_output.power.counter = global_summary.result.power.value.unwrap_or(0);
        jeedom_summary_output.security.counter = global_summary.result.security.value.unwrap_or(0);
        jeedom_summary_output.shutter.counter = global_summary.result.shutter.value.unwrap_or(0);
        jeedom_summary_output.temperature.counter =
            global_summary.result.temperature.value.unwrap_or(0);
        jeedom_summary_output.windows.counter = global_summary.result.windows.value.unwrap_or(0);

        // update notifications
        jeedom_summary_output.update_notification_info(notifications);

        // update battery status
        jeedom_summary_output.update_battery_info(battery_status);

        jeedom_summary_output
    }

    // icons: https://www.nerdfonts.com/cheat-sheet
    fn init_with_nerd_fonts_icons() -> Self {
        Self {
            alarm: JeedomSummaryItemOutput::new("\u{F023}".to_string()),
            battery: JeedomBatteryInfoOuput::new("\u{F244}".to_string()),
            door: JeedomSummaryItemOutput::new("\u{FD18}".to_string()), // not ok
            humidity: JeedomSummaryItemOutput::new("\u{E373}".to_string()),
            light: JeedomSummaryItemOutput::new("\u{F834}".to_string()),
            luminosity: JeedomSummaryItemOutput::new("\u{FAA7}".to_string()),
            motion: JeedomSummaryItemOutput::new("\u{FC0C}".to_string()), // not ok
            outlet: JeedomSummaryItemOutput::new("\u{F1E6}".to_string()),
            power: JeedomSummaryItemOutput::new("\u{F0E7}".to_string()),
            security: JeedomSummaryItemOutput::new("\u{FC8D}".to_string()), //not ok
            shutter: JeedomSummaryItemOutput::new("S".to_string()),
            temperature: JeedomSummaryItemOutput::new("\u{F2C7}".to_string()),
            updates: JeedomSummaryItemOutput::new("\u{F62E}".to_string()),
            windows: JeedomSummaryItemOutput::new("\u{F17A}".to_string()),
            notifications: JeedomNotificationsItemOutput::default(),
        }
    }

    // emoji: https://unicode.org/emoji/charts/full-emoji-list.html
    fn init_with_emoji_icons() -> Self {
        Self {
            alarm: JeedomSummaryItemOutput::new("\u{1F512}".to_string()),
            battery: JeedomBatteryInfoOuput::new("\u{1F50}".to_string()),
            door: JeedomSummaryItemOutput::new("\u{1F6AA}".to_string()),
            humidity: JeedomSummaryItemOutput::new("\u{1F4A7}".to_string()),
            light: JeedomSummaryItemOutput::new("\u{1F4A1}".to_string()),
            luminosity: JeedomSummaryItemOutput::new("\u{1F506}".to_string()),
            motion: JeedomSummaryItemOutput::new("\u{1F3C3}".to_string()),
            outlet: JeedomSummaryItemOutput::new("\u{1F50C}".to_string()),
            power: JeedomSummaryItemOutput::new("\u{26A1}".to_string()),
            security: JeedomSummaryItemOutput::new("\u{1F6A8}".to_string()),
            shutter: JeedomSummaryItemOutput::new("\u{2195}".to_string()),
            temperature: JeedomSummaryItemOutput::new("\u{1F321}".to_string()),
            updates: JeedomSummaryItemOutput::new("\u{1F534}".to_string()),
            windows: JeedomSummaryItemOutput::new("\u{1F5BC}".to_string()),
            notifications: JeedomNotificationsItemOutput::default(),
        }
    }

    // Load fonts with http://mathew-kurian.github.io/CharacterMap/
    fn init_with_jeedom_icons() -> Self {
        Self {
            alarm: JeedomSummaryItemOutput::new("\u{E60E}".to_string()), //Jeedom font
            battery: JeedomBatteryInfoOuput::new("\u{E602}".to_string()), //Jeedom font
            door: JeedomSummaryItemOutput::new("\u{E61D}".to_string()),  //Jeedom font
            humidity: JeedomSummaryItemOutput::new("\u{E90F}".to_string()), //Jeedomapp font
            light: JeedomSummaryItemOutput::new("\u{E611}".to_string()), //Jeedom font
            luminosity: JeedomSummaryItemOutput::new("\u{E601}".to_string()), //Nature font
            motion: JeedomSummaryItemOutput::new("\u{E612}".to_string()), //Jeedom font
            outlet: JeedomSummaryItemOutput::new("\u{E61E}".to_string()), //Jeedom font
            power: JeedomSummaryItemOutput::new("\u{F0E7}".to_string()), //General font / fonts awesome
            security: JeedomSummaryItemOutput::new("\u{E601}".to_string()), //Jeedom font
            shutter: JeedomSummaryItemOutput::new("\u{E627}".to_string()), //Jeedom font
            temperature: JeedomSummaryItemOutput::new("\u{E622}".to_string()), //Jeedom font
            updates: JeedomSummaryItemOutput::new("\u{E91D}".to_string()), // Jeedomapp font
            windows: JeedomSummaryItemOutput::new("\u{E60A}".to_string()), //Jeedom font
            notifications: JeedomNotificationsItemOutput::default(),
        }
    }

    fn init_without_icons() -> Self {
        Self {
            alarm: JeedomSummaryItemOutput::new("A".to_string()),
            battery: JeedomBatteryInfoOuput::new("B".to_string()),
            door: JeedomSummaryItemOutput::new("D".to_string()),
            humidity: JeedomSummaryItemOutput::new("H".to_string()),
            light: JeedomSummaryItemOutput::new("G".to_string()),
            luminosity: JeedomSummaryItemOutput::new("L".to_string()),
            motion: JeedomSummaryItemOutput::new("M".to_string()),
            outlet: JeedomSummaryItemOutput::new("O".to_string()),
            power: JeedomSummaryItemOutput::new("P".to_string()),
            security: JeedomSummaryItemOutput::new("S".to_string()),
            shutter: JeedomSummaryItemOutput::new("U".to_string()),
            temperature: JeedomSummaryItemOutput::new("R".to_string()),
            updates: JeedomSummaryItemOutput::new("U".to_string()),
            windows: JeedomSummaryItemOutput::new("W".to_string()),
            notifications: JeedomNotificationsItemOutput::default(),
        }
    }

    pub fn is_there_something_to_print(&self, ignore_battery_warning: bool) -> bool {
        if self.alarm.counter > 0
            || self.battery.battery_danger_counter > 0
            || (ignore_battery_warning && self.battery.battery_warning_counter > 0)
            || self.door.counter > 0
            || self.humidity.counter > 0
            || self.light.counter > 0
            || self.luminosity.counter > 0
            || self.motion.counter > 0
            || self.outlet.counter > 0
            || self.power.counter > 0
            || self.security.counter > 0
            || self.shutter.counter > 0
            || self.temperature.counter > 0
            || self.updates.counter > 0
            || self.windows.counter > 0
            || self.notifications.message_counter > 0
            || self.notifications.update_counter > 0
        {
            return true;
        }
        false
    }

    pub fn update_battery_info(&mut self, battery_status: JeedomBatteryResult) {
        for battery in battery_status.result {
            if let Some(x) = battery.status {
                let battery_level = match x.battery {
                    Some(level) => level,
                    None => BatteryPercentage(Some(100)), // set to 100 to check warning and danger anyway and ignore level
                };
                // danger
                match x.batterydanger {
                    Some(danger) if danger == 1 => {
                        self.battery.battery_danger_counter += 1;
                        continue; // because warning = 1 when danger = 1
                    }
                    Some(_) => {
                        if let Some(battery_level) = battery_level.0 {
                            if battery_level <= JEEDOM_BATTERY_MIN_DANGER {
                                self.battery.battery_danger_counter += 1;
                                continue; // because warning = 1 when danger = 1
                            }
                        }
                    }
                    None => {}
                }
                // warning
                match x.batterywarning {
                    Some(warning) if warning == 1 => {
                        self.battery.battery_warning_counter += 1;
                    }
                    Some(_) => {
                        if let Some(battery_level) = battery_level.0 {
                            if battery_level <= JEEDOM_BATTERY_MIN_WARNING {
                                self.battery.battery_warning_counter += 1;
                            }
                        }
                    }
                    None => {}
                }
            }
        }
    }

    pub fn update_notification_info(&mut self, notification_status: JeedomNotificationResult) {
        // notifications messages
        self.notifications.message_counter = notification_status.result.len() as u32;
        self.notifications.message_icon = match JEEDOM_NOTIFICATION_ICONS
            .get(format!("{}", self.notifications.message_counter).as_str())
        {
            Some(x) => Some(x.to_string()),
            None => Some(self.notifications.message_counter.to_string()),
        };

        // update messages
        self.notifications.update_counter = notification_status
            .result
            .iter()
            .filter(|x| x.plugin == "update")
            .count() as u32;
        self.notifications.update_icon = match JEEDOM_NOTIFICATION_ICONS
            .get(format!("{}", self.notifications.update_counter).as_str())
        {
            Some(x) => Some(x.to_string()),
            None => Some(self.notifications.update_counter.to_string()),
        };
    }

    pub fn render_bar_output(
        &self,
        bar_type: BarType,
        ignore_battery_warning: bool,
        jeedom_url: String,
    ) -> String {
        if !self.is_there_something_to_print(ignore_battery_warning) {
            return "Jeedom".to_string();
        }
        let mut jeedom_summary_bar_content = Vec::new();

        // first part: jeedom global summary
        if self.alarm.counter > 0 {
            jeedom_summary_bar_content.push(format!("{}{}", self.alarm.counter, self.alarm.icon));
        }
        if self.door.counter > 0 {
            jeedom_summary_bar_content.push(format!("{}{}", self.door.counter, self.door.icon));
        }
        if self.humidity.counter > 0 {
            jeedom_summary_bar_content
                .push(format!("{}{}", self.humidity.counter, self.humidity.icon));
        }
        if self.light.counter > 0 {
            jeedom_summary_bar_content.push(format!("{}{}", self.light.counter, self.light.icon));
        }
        if self.luminosity.counter > 0 {
            jeedom_summary_bar_content.push(format!(
                "{}{}",
                self.luminosity.counter, self.luminosity.icon
            ));
        }
        if self.motion.counter > 0 {
            jeedom_summary_bar_content.push(format!("{}{}", self.motion.counter, self.motion.icon));
        }
        if self.outlet.counter > 0 {
            jeedom_summary_bar_content.push(format!("{}{}", self.outlet.counter, self.outlet.icon));
        }
        if self.power.counter > 0 {
            jeedom_summary_bar_content.push(format!("{}{}", self.power.counter, self.power.icon));
        }
        if self.security.counter > 0 {
            jeedom_summary_bar_content
                .push(format!("{}{}", self.security.counter, self.security.icon));
        }
        if self.shutter.counter > 0 {
            jeedom_summary_bar_content
                .push(format!("{}{}", self.shutter.counter, self.shutter.icon));
        }
        if self.temperature.counter > 0 {
            jeedom_summary_bar_content.push(format!(
                "{}{}",
                self.temperature.counter, self.temperature.icon
            ));
        }
        if self.windows.counter > 0 {
            jeedom_summary_bar_content
                .push(format!("{}{}", self.windows.counter, self.windows.icon));
        }

        // notifications
        match bar_type {
            BarType::Mac => {
                if self.notifications.update_counter > 0 {
                    if let Some(x) = &self.notifications.update_icon {
                        jeedom_summary_bar_content.push(x.red().to_string());
                    }
                };
                if self.notifications.message_counter > 0 {
                    if let Some(x) = &self.notifications.message_icon {
                        jeedom_summary_bar_content.push(x.yellow().to_string());
                    }
                };
            }
            BarType::I3blocks | BarType::I3StatusRust => {
                if self.notifications.update_counter > 0 {
                    if let Some(x) = &self.notifications.update_icon {
                        jeedom_summary_bar_content.push(format!(
                            "<span color='red'><span font='Jeedom'>{x}</span></span>",
                        ));
                    }
                };
                if self.notifications.message_counter > 0 {
                    if let Some(x) = &self.notifications.message_icon {
                        jeedom_summary_bar_content.push(format!(
                            "<span color='yellow'><span font='Jeedom'>{x}</span></span>",
                        ));
                    }
                };
            }
            BarType::None => {
                if self.notifications.update_counter > 0 {
                    if let Some(x) = &self.notifications.update_icon {
                        jeedom_summary_bar_content.push(x.to_string());
                    }
                };
                if self.notifications.message_counter > 0 {
                    if let Some(x) = &self.notifications.message_icon {
                        jeedom_summary_bar_content.push(x.to_string());
                    }
                };
            }
        };

        // last part: battery
        match bar_type {
            BarType::Mac => {
                // warning battery
                if !ignore_battery_warning && self.battery.battery_warning_counter > 0 {
                    jeedom_summary_bar_content.push(
                        format!(
                            "{}{}",
                            self.battery.battery_warning_counter, self.battery.icon
                        )
                        .yellow()
                        .to_string(),
                    );
                };
                // danger battery
                if self.battery.battery_danger_counter > 0 {
                    jeedom_summary_bar_content.push(
                        format!(
                            "{}{}",
                            self.battery.battery_danger_counter, self.battery.icon
                        )
                        .red()
                        .to_string(),
                    );
                }
            }
            BarType::I3blocks | BarType::I3StatusRust => {
                // warning battery
                if !ignore_battery_warning && self.battery.battery_warning_counter > 0 {
                    jeedom_summary_bar_content.push(format!(
                        "<span color='yellow'><span font='Jeedom'>{}{}</span></span>",
                        self.battery.battery_warning_counter, self.battery.icon
                    ));
                };
                // danger battery
                if self.battery.battery_danger_counter > 0 {
                    jeedom_summary_bar_content.push(format!(
                        "<span color='red'><span font='Jeedom'>{}{}</span></span>",
                        self.battery.battery_danger_counter, self.battery.icon
                    ));
                }
            }
            BarType::None => {
                // warning battery
                if !ignore_battery_warning && self.battery.battery_warning_counter > 0 {
                    jeedom_summary_bar_content.push(format!(
                        "{}{}",
                        self.battery.battery_warning_counter, self.battery.icon
                    ));
                };
                // danger battery
                if self.battery.battery_danger_counter > 0 {
                    jeedom_summary_bar_content.push(format!(
                        "{}{}",
                        self.battery.battery_danger_counter, self.battery.icon
                    ));
                }
            }
        };

        let jeedom_summary_bar_final = jeedom_summary_bar_content.join(" ");
        match bar_type {
            BarType::Mac => {
                let mut mac_output = vec!["---".to_string()];
                if self.notifications.update_counter > 0 {
                    mac_output.push(format!(
                        "Updates {} | color=red href={}/index.php?v=d&p=update",
                        self.notifications.update_counter, jeedom_url
                    ))
                };
                if self.notifications.message_counter > 0 {
                    mac_output.push(format!(
                        "Messages {} | color=orange href={}",
                        self.notifications.message_counter, jeedom_url
                    ))
                };
                mac_output.join("\n")
            }
            BarType::I3blocks => format!("{jeedom_summary_bar_final}\n{jeedom_summary_bar_final}"),
            BarType::I3StatusRust => jeedom_summary_bar_final,
            BarType::None => format!("{jeedom_summary_bar_final}\n"),
        }
    }
}

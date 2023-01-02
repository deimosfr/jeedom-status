use crate::api::JeedomClient;
use clap::{Parser, ValueEnum};
use output::JeedomSummaryOutput;

pub mod api;
pub mod io_models;
pub mod output;

// Get Jeedom global summary
#[derive(Parser, Debug)]
#[command(about = "Get summary of your Jeedom equipments, notifications and battery", long_about = None)]
#[command(version = "2.0.0")] // ci-version-check
struct Args {
    /// Jeedom Local API URL, like http://jeedom (required)
    #[arg(short = 'u', long = "url", env = "JEEDOM_URL")]
    url: String,
    /// Jeedom Extrenal API URL, like http://jeedom
    #[arg(short = 'a', long = "alertnateUrl")]
    alertnate_url: Option<String>,
    /// Jeedom API key or User Hash Key (required)
    #[arg(short = 'k', long = "apiKey", env = "JEEDOM_API_KEY")]
    api_key: String,
    // Select the bar type
    #[arg(short = 'b', long = "barType", value_enum, default_value = "mac")]
    bar_type: BarType,
    // Select the bar style
    #[arg(short = 's', long = "barStyle", value_enum, default_value = "text")]
    bar_style: BarStyle,
    /// Ignore battery warning report
    #[arg(short = 'w', long = "ignore-battery-warning", default_value = "false")]
    ignore_battery_warning: bool,
    /// Run in debug mode
    #[arg(short = 'd', long = "debug", default_value = "false", env = "DEBUG")]
    enable_debug_mode: bool,
}

#[derive(Copy, Clone, PartialEq, Eq, PartialOrd, Ord, ValueEnum, Debug)]
pub enum BarType {
    Mac,
    I3blocks,
    I3StatusRust,
    None,
}

#[derive(Copy, Clone, PartialEq, Eq, PartialOrd, Ord, ValueEnum, Debug)]
pub enum BarStyle {
    Text,
    Jeedom,
    Nerd,
    Emoji,
}

fn main() {
    let args = Args::parse();

    if args.enable_debug_mode {
        println!("Used args:\n{args:?}");
    }

    // automatically select the reachable Jeedom URL
    let url_to_use = match JeedomClient::check_connectivity(
        args.url,
        args.alertnate_url,
        args.enable_debug_mode,
    ) {
        Ok(x) => x,
        Err(e) => {
            println!("Jeedom N/A");
            if args.enable_debug_mode {
                println!("{e}");
            }
            std::process::exit(1);
        }
    };

    // client
    let jeedom_client = JeedomClient::new(url_to_use, args.api_key, args.enable_debug_mode);

    // ensure Jeedom is reachable (required behind proxy)
    if let Err(e) = jeedom_client.ping() {
        if args.enable_debug_mode {
            println!("Jeedom ping fail: {e}");
        }
        println!("Jeedom N/A");
        std::process::exit(1);
    }

    // get global summary
    let global_summary = match jeedom_client.global_summary() {
        Ok(x) => x,
        Err(e) => {
            println!("Global summary Jeedom error: {e}");
            std::process::exit(1);
        }
    };

    // get all battery equiments info
    let battery_status = match jeedom_client.battery_status() {
        Ok(x) => x,
        Err(e) => {
            println!("Battery status Jeedom error: {e}");
            std::process::exit(1);
        }
    };

    // get all notifications
    let notifications = match jeedom_client.notification_messages() {
        Ok(x) => x,
        Err(e) => {
            println!("Notifications Jeedom error: {e}");
            std::process::exit(1);
        }
    };

    // build output
    let jeedom_summary_output = JeedomSummaryOutput::new(
        global_summary,
        battery_status,
        notifications,
        args.bar_style,
    );

    // print bar
    let bar_output = jeedom_summary_output.render_bar_output(
        args.bar_type,
        args.ignore_battery_warning,
        jeedom_client.url,
    );
    println!("{bar_output}");
}

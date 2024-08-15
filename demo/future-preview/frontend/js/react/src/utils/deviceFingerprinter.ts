export async function getBrowserInfoWithHash() {
  const ua = navigator.userAgent;
  let browserName = "Unknown";
  let browserVersion = "Unknown";
  let osName = "Unknown";

  // Detect browser name
  if (ua.indexOf("Firefox") > -1) {
    browserName = "Firefox";
  } else if (ua.indexOf("Opera") > -1 || ua.indexOf("OPR") > -1) {
    browserName = "Opera";
  } else if (ua.indexOf("Trident") > -1) {
    browserName = "Internet Explorer";
  } else if (ua.indexOf("Edge") > -1) {
    browserName = "Edge";
  } else if (ua.indexOf("Chrome") > -1) {
    browserName = "Chrome";
  } else if (ua.indexOf("Safari") > -1) {
    browserName = "Safari";
  }

  // Detect browser version
  const match =
    ua.match(/(opera|chrome|safari|firefox|msie|trident(?=\/))\/?\s*(\d+)/i) ||
    [];
  if (match.length > 2) {
    browserVersion = match[2];
  }

  // Detect OS
  if (navigator.appVersion.indexOf("Win") != -1) osName = "Windows";
  if (navigator.appVersion.indexOf("Mac") != -1) osName = "MacOS";
  if (navigator.appVersion.indexOf("X11") != -1) osName = "UNIX";
  if (navigator.appVersion.indexOf("Linux") != -1) osName = "Linux";

  const browserInfo = {
    browser: browserName,
    browserVersion: browserVersion,
    os: osName,
    userAgent: ua,
    language: navigator.language,
    platform: navigator.platform,
    cookiesEnabled: navigator.cookieEnabled,
    screenResolution: `${window.screen.width}x${window.screen.height}`,
    colorDepth: window.screen.colorDepth,
    referrer: document.referrer,
    timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  };

  // Generate hash using Web Crypto API
  const infoString = JSON.stringify(browserInfo);
  const encoder = new TextEncoder();
  const data = encoder.encode(infoString);
  const hashBuffer = await crypto.subtle.digest("SHA-256", data);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
}

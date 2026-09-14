interface ImportMetaEnv {
  readonly VITE_VEHICLE_SERVICE_URL: string;
  readonly VITE_ALERT_SERVICE_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

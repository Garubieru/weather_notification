export type WeatherNotificationContent = {
  cityName: string;
  cityStateCode: string;
  prediction: {
    temperatures: Array<{
      date: string;
      min: number;
      max: number;
      condition: string;
      uvi: number;
    }>;
    waveConditions: {
      morning: string;
      afternoon: string;
      evening: string;
    };
  };
};

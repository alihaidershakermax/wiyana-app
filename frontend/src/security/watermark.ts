export interface WatermarkConfig {
  displayName: string;
  telegramId: string;
  privacyId: string;
  sessionId: string;
}

export interface WatermarkPosition {
  topPercent: number;
  leftPercent: number;
  opacity: number;
  rotationDeg: number;
}

export function generateRandomPosition(containerWidth = 100, containerHeight = 100): WatermarkPosition {
  // Safe boundaries (10% to 75% for top, 10% to 70% for left)
  const topPercent = Math.floor(Math.random() * 65) + 10;
  const leftPercent = Math.floor(Math.random() * 60) + 10;
  const opacity = Math.round((Math.random() * 0.35 + 0.35) * 100) / 100; // 0.35 to 0.70
  const rotationDeg = Math.floor(Math.random() * 14) - 7; // -7deg to +7deg

  return { topPercent, leftPercent, opacity, rotationDeg };
}

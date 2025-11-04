// Pastel color palette for modern, minimal design
export const colors = {
  // Primary pastel colors
  primary: '#A8D5E2', // Pastel blue
  secondary: '#F9A8D4', // Pastel pink
  accent: '#C4B5FD', // Pastel purple
  success: '#BBF7D0', // Pastel green
  warning: '#FED7AA', // Pastel orange
  error: '#FECACA', // Pastel red

  // Neutral colors
  background: '#FEFEFE',
  surface: '#F8F9FA',
  surfaceLight: '#FFFFFF',
  border: '#E5E7EB',
  borderLight: '#F3F4F6',

  // Text colors
  text: '#1F2937',
  textSecondary: '#6B7280',
  textLight: '#9CA3AF',
  textInverse: '#FFFFFF',

  // Streak colors
  streakActive: '#BBF7D0',
  streakInactive: '#F3F4F6',
  streakToday: '#A8D5E2',
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
};

export const borderRadius = {
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  full: 9999,
};

export const typography = {
  // Font sizes
  xs: 12,
  sm: 14,
  base: 16,
  lg: 18,
  xl: 20,
  xxl: 24,
  xxxl: 32,

  // Font weights
  regular: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: '700' as const,
};

export const shadows = {
  sm: {
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 1},
    shadowOpacity: 0.05,
    shadowRadius: 2,
    elevation: 2,
  },
  md: {
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.08,
    shadowRadius: 4,
    elevation: 3,
  },
  lg: {
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 4},
    shadowOpacity: 0.1,
    shadowRadius: 8,
    elevation: 5,
  },
};

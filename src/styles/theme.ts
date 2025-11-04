/**
 * Theme and styling constants for the app
 * Modern minimal design with pastel colors and rounded corners
 */

export const colors = {
  // Pastel color palette
  primary: '#A8DADC', // Pastel cyan
  secondary: '#F1FAEE', // Pastel cream
  accent: '#E63946', // Muted red for important actions
  success: '#B8E6B8', // Pastel green
  warning: '#FFE5B4', // Pastel peach
  error: '#FFB3BA', // Pastel red

  // Neutral colors
  background: '#FEFEFE',
  surface: '#F8F9FA',
  surfaceLight: '#FFFFFF',
  text: '#2D3436',
  textLight: '#636E72',
  textMuted: '#95A5A6',
  border: '#DFE6E9',
  borderLight: '#E8EEF1',

  // Streak colors
  streakActive: '#A8E6CF', // Pastel mint green
  streakInactive: '#E8EEF1', // Very light gray
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
  round: 999,
};

export const fontSize = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 22,
  xxl: 28,
  xxxl: 36,
};

export const fontFamily = {
  regular: 'System', // Uses system sans-serif font
  medium: 'System',
  bold: 'System',
};

export const shadow = {
  sm: {
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 1},
    shadowOpacity: 0.1,
    shadowRadius: 2,
    elevation: 2,
  },
  md: {
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.15,
    shadowRadius: 4,
    elevation: 4,
  },
  lg: {
    shadowColor: '#000',
    shadowOffset: {width: 0, height: 4},
    shadowOpacity: 0.2,
    shadowRadius: 8,
    elevation: 8,
  },
};

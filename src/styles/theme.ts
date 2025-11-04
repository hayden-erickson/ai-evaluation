/**
 * Theme configuration with pastel colors and modern design
 */

export const theme = {
  colors: {
    // Pastel primary colors
    primary: '#B4A5E8', // Soft lavender
    primaryLight: '#D4C9F2',
    primaryDark: '#9687D4',

    // Pastel secondary colors
    secondary: '#A8D8EA', // Soft blue
    secondaryLight: '#C8E8F4',
    secondaryDark: '#88B8CA',

    // Pastel accent colors
    accent: '#FFDAC1', // Soft peach
    accentLight: '#FFE9D6',
    accentDark: '#FFCAA6',

    success: '#B8E6C9', // Soft green
    warning: '#FFE4B5', // Soft yellow
    error: '#FFB8B8', // Soft red
    info: '#B8D8F0', // Soft sky blue

    // Neutral colors
    background: '#FAFAFA',
    surface: '#FFFFFF',
    surfaceLight: '#F5F5F5',
    
    text: '#333333',
    textLight: '#666666',
    textLighter: '#999999',
    
    border: '#E0E0E0',
    borderLight: '#F0F0F0',

    // Streak colors
    streakActive: '#B8E6C9',
    streakEmpty: '#F0F0F0',
  },
  
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 40,
  },
  
  borderRadius: {
    sm: 8,
    md: 12,
    lg: 16,
    xl: 20,
    round: 50,
  },
  
  fontSize: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 18,
    xl: 24,
    xxl: 32,
  },
  
  fontFamily: {
    regular: 'System',
    medium: 'System',
    bold: 'System',
  },
  
  shadows: {
    small: {
      shadowColor: '#000',
      shadowOffset: {width: 0, height: 2},
      shadowOpacity: 0.05,
      shadowRadius: 4,
      elevation: 2,
    },
    medium: {
      shadowColor: '#000',
      shadowOffset: {width: 0, height: 4},
      shadowOpacity: 0.08,
      shadowRadius: 8,
      elevation: 4,
    },
    large: {
      shadowColor: '#000',
      shadowOffset: {width: 0, height: 8},
      shadowOpacity: 0.12,
      shadowRadius: 16,
      elevation: 8,
    },
  },
};

export type Theme = typeof theme;

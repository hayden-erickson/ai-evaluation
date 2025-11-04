/**
 * App theme with pastel colors and design tokens
 */

export const theme = {
  colors: {
    // Pastel color palette
    primary: '#B4A7D6',        // Soft lavender
    secondary: '#A7C7D6',      // Soft blue
    accent: '#D6B4C2',         // Soft pink
    success: '#B4D6B4',        // Soft green
    warning: '#D6D4B4',        // Soft yellow
    error: '#D6B4B4',          // Soft red
    
    // Neutral colors
    background: '#FAFAFA',     // Very light gray
    surface: '#FFFFFF',        // White
    text: '#4A4A4A',          // Dark gray
    textSecondary: '#8A8A8A', // Medium gray
    border: '#E0E0E0',        // Light gray
    
    // Overlay
    overlay: 'rgba(0, 0, 0, 0.5)',
  },
  
  spacing: {
    xs: 4,
    sm: 8,
    md: 16,
    lg: 24,
    xl: 32,
    xxl: 48,
  },
  
  borderRadius: {
    sm: 8,
    md: 12,
    lg: 16,
    xl: 24,
    round: 999,
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
    sm: {
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 1 },
      shadowOpacity: 0.05,
      shadowRadius: 2,
      elevation: 1,
    },
    md: {
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 2 },
      shadowOpacity: 0.08,
      shadowRadius: 4,
      elevation: 3,
    },
    lg: {
      shadowColor: '#000',
      shadowOffset: { width: 0, height: 4 },
      shadowOpacity: 0.12,
      shadowRadius: 8,
      elevation: 5,
    },
  },
};

export type Theme = typeof theme;


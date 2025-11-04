/**
 * Global styles and theme configuration
 * Modern, minimal design with pastel colors and rounded corners
 */

import { StyleSheet } from 'react-native';

// Color palette - Pastel colors for modern, minimal design
export const Colors = {
  // Primary pastel colors
  primary: '#A8DADC',      // Soft cyan
  secondary: '#F1FAEE',    // Off white
  accent: '#E63946',       // Soft red
  success: '#B8E6D5',      // Soft mint green
  warning: '#FFD6A5',      // Soft peach
  
  // Background colors
  background: '#F8F9FA',   // Light gray
  cardBackground: '#FFFFFF',
  modalBackground: 'rgba(0, 0, 0, 0.5)',
  
  // Text colors
  textPrimary: '#2B2D42',  // Dark blue-gray
  textSecondary: '#8D99AE', // Medium gray
  textLight: '#FFFFFF',
  
  // UI element colors
  border: '#E0E0E0',
  inputBackground: '#F5F5F5',
  disabled: '#D3D3D3',
  
  // Streak colors
  streakActive: '#B8E6D5',
  streakInactive: '#E0E0E0',
};

// Typography - Sans serif fonts
export const Typography = {
  fontFamily: 'System',  // Uses system default sans-serif
  
  // Font sizes
  h1: 32,
  h2: 24,
  h3: 20,
  h4: 18,
  body: 16,
  small: 14,
  tiny: 12,
  
  // Font weights
  regular: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: '700' as const,
};

// Spacing
export const Spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
};

// Border radius for rounded corners
export const BorderRadius = {
  sm: 8,
  md: 12,
  lg: 16,
  xl: 24,
  full: 9999,
};

// Common styles used across components
export const CommonStyles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.background,
  },
  
  card: {
    backgroundColor: Colors.cardBackground,
    borderRadius: BorderRadius.md,
    padding: Spacing.md,
    marginBottom: Spacing.md,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 2,
  },
  
  button: {
    backgroundColor: Colors.primary,
    borderRadius: BorderRadius.md,
    paddingVertical: Spacing.md,
    paddingHorizontal: Spacing.lg,
    alignItems: 'center',
    justifyContent: 'center',
  },
  
  buttonText: {
    color: Colors.textPrimary,
    fontSize: Typography.body,
    fontWeight: Typography.semibold,
  },
  
  input: {
    backgroundColor: Colors.inputBackground,
    borderRadius: BorderRadius.sm,
    padding: Spacing.md,
    fontSize: Typography.body,
    color: Colors.textPrimary,
    borderWidth: 1,
    borderColor: Colors.border,
  },
  
  textPrimary: {
    color: Colors.textPrimary,
    fontSize: Typography.body,
    fontWeight: Typography.regular,
  },
  
  textSecondary: {
    color: Colors.textSecondary,
    fontSize: Typography.small,
    fontWeight: Typography.regular,
  },
  
  heading1: {
    color: Colors.textPrimary,
    fontSize: Typography.h1,
    fontWeight: Typography.bold,
  },
  
  heading2: {
    color: Colors.textPrimary,
    fontSize: Typography.h2,
    fontWeight: Typography.bold,
  },
  
  heading3: {
    color: Colors.textPrimary,
    fontSize: Typography.h3,
    fontWeight: Typography.semibold,
  },
  
  modalContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: Colors.modalBackground,
  },
  
  modalContent: {
    backgroundColor: Colors.cardBackground,
    borderRadius: BorderRadius.lg,
    padding: Spacing.xl,
    width: '90%',
    maxWidth: 400,
  },
  
  errorText: {
    color: Colors.accent,
    fontSize: Typography.small,
    marginTop: Spacing.xs,
  },
  
  centered: {
    alignItems: 'center',
    justifyContent: 'center',
  },
  
  row: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  
  spaceBetween: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
  },
});

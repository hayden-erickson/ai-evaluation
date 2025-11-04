import { createTheme } from '@rneui/themed';
import { colors } from './colors';

export const theme = createTheme({
  lightColors: {
    primary: colors.primary,
    secondary: colors.secondary,
    background: colors.background,
    error: colors.error,
  },
  components: {
    Button: {
      radius: 8,
      buttonStyle: {
        backgroundColor: colors.primary,
      },
      titleStyle: {
        fontFamily: 'System', // sans-serif
      },
    },
    Input: {
      inputContainerStyle: {
        borderBottomWidth: 1,
        borderColor: colors.gray,
      },
      inputStyle: {
        fontFamily: 'System',
      },
    },
    Text: {
      style: {
        fontFamily: 'System',
        color: colors.text,
      },
    },
  },
});

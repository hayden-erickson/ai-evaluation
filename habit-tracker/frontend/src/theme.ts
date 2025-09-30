import { extendTheme } from '@chakra-ui/theme';

const theme = extendTheme({
    styles: {
        global: {
            body: {
                bg: 'gray.50',
            },
        },
    },
    components: {
        Button: {
            defaultProps: {
                colorScheme: 'teal',
            },
        },
    },
});

export default theme;
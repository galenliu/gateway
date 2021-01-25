import {createGlobalState} from 'react-hooks-global-state';

export const {useGlobalState} = createGlobalState({

    menuState: {
        menuBtuShow: true,
        menuListShow: false,
        menuScrimShow: false,
    }
});

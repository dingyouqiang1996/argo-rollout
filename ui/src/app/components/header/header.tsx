import * as React from 'react';

import {InfoItemRow, Theme, ThemeContext, ThemeDiv} from 'argo-ui/v2';
import {useParams} from 'react-router';
import {NamespaceContext, RolloutAPIContext} from '../../shared/context/api';

import './header.scss';
import {Link, useHistory} from 'react-router-dom';
import {AutoComplete, Button, Tooltip} from 'antd';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faBook, faKeyboard, faMoon, faSun} from '@fortawesome/free-solid-svg-icons';

export const ThemeToggle = () => {
    const dmCtx = React.useContext(ThemeContext);
    const isDark = dmCtx.theme === Theme.Dark;
    const icon = isDark ? faSun : faMoon;
    return <Button onClick={() => dmCtx.set(isDark ? Theme.Light : Theme.Dark)} icon={<FontAwesomeIcon icon={icon} />} />;
};

const Logo = () => <img src='assets/images/argo-icon-color-square.png' style={{width: '37px', height: '37px', margin: '0 12px'}} alt='Argo Logo' />;

export const Header = (props: {pageHasShortcuts: boolean; changeNamespace: (val: string) => void; showHelp: () => void}) => {
    const history = useHistory();
    const namespaceInfo = React.useContext(NamespaceContext);
    const {namespace} = useParams<{namespace: string}>();
    const api = React.useContext(RolloutAPIContext);
    const [version, setVersion] = React.useState('v?');
    const [nsInput, setNsInput] = React.useState(namespaceInfo.namespace);
    React.useEffect(() => {
        const getVersion = async () => {
            const v = await api.rolloutServiceVersion();
            setVersion(v.rolloutsVersion);
        };
        getVersion();
    }, []);
    React.useEffect(() => {
        if (namespace && namespace != namespaceInfo.namespace) {
            props.changeNamespace(namespace);
            setNsInput(namespace);
        }
    }, []);
    return (
        <header className='rollouts-header'>
            <Link to='/' className='rollouts-header__brand'>
                <Logo />
                <div>
                    <div className='rollouts-header__title'>
                        <img src='/assets/images/argologo.svg' alt='Argo Text Logo' style={{filter: 'invert(100%)', height: '1em'}}/>
                    </div>
                    <div className='rollouts-header__label'>Rollouts {version}</div>
                </div>
            </Link>
            <div className='rollouts-header__info'>
                {props.pageHasShortcuts && (
                    <Tooltip title='Keyboard Shortcuts'>
                        <Button onClick={props.showHelp} icon={<FontAwesomeIcon icon={faKeyboard} />} style={{marginRight: '10px'}} />
                    </Tooltip>
                )}
                <Tooltip title='Documentation'>
                    <a href='https://argo-rollouts.readthedocs.io/' target='_blank' rel='noreferrer'>
                        <Button icon={<FontAwesomeIcon icon={faBook} />} style={{marginRight: '10px'}} />
                    </a>
                </Tooltip>
                <span style={{marginRight: '7px'}}>
                    <Tooltip title='Toggle Dark Mode'>
                        <ThemeToggle />
                    </Tooltip>
                </span>
                {(namespaceInfo.availableNamespaces || []).length == 0 ? (
                    <InfoItemRow label={'NS:'} items={{content: namespaceInfo.namespace}} />
                ) : (
                    <ThemeDiv className='rollouts-header__namespace'>
                        <div className='rollouts-header__label'>NAMESPACE</div>
                        <AutoComplete
                            style={{width: 200}}
                            className='rollouts-header__namespace-selector'
                            options={(namespaceInfo.availableNamespaces || []).map((ns) => ({label: ns, value: ns}))}
                            placeholder='Namespace'
                            onChange={(val) => setNsInput(val)}
                            onSelect={(val) => {
                                const selectedNamespace = val ? val : nsInput;
                                props.changeNamespace(selectedNamespace);
                                history.push(`/${selectedNamespace}`);
                            }}
                            value={nsInput}
                        />
                    </ThemeDiv>
                )}
            </div>
        </header>
    );
};

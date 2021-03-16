import * as React from 'react';
import {useParams} from 'react-router-dom';
import {Helmet} from 'react-helmet';

import './rollout.scss';
import {RolloutStatus, StatusIcon} from '../status-icon/status-icon';
import {ThemeDiv} from '../theme-div/theme-div';
import {useWatchRollout} from '../../shared/services/rollout';
import {InfoItem, InfoItemKind, InfoItemProps, InfoItemRow} from '../info-item/info-item';
import {RolloutInfo} from '../../../models/rollout/rollout';
import {
    faBalanceScale,
    faBalanceScaleRight,
    faCheck,
    faChevronCircleDown,
    faChevronCircleUp,
    faClock,
    faDove,
    faPalette,
    faPauseCircle,
    faPencilAlt,
    faShoePrints,
    faTimes,
    faWeight,
    IconDefinition,
} from '@fortawesome/free-solid-svg-icons';
import {ReplicaSet} from '../pods/pods';
import {formatTimestamp, IconForTag, ImageTag} from '../../shared/utils/utils';
import {RolloutAPIContext} from '../../shared/context/api';
import {useInput} from '../input/input';
import {ActionButton} from '../action-button/action-button';
import {Spinner, WaitFor} from '../wait-for/wait-for';
import {
    GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1AnalysisRunInfo,
    GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1CanaryStep,
    GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1ContainerInfo,
    GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1ExperimentInfo,
    GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1ReplicaSetInfo,
} from '../../../models/rollout/generated';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {Autocomplete} from '../autocomplete/autocomplete';
import {faChartBar} from '@fortawesome/free-regular-svg-icons';
const RolloutActions = React.lazy(() => import('../rollout-actions/rollout-actions'));
interface ImageInfo {
    image: string;
    tags: ImageTag[];
    color?: ImageColor;
}

enum ImageColor {
    Red = 'red',
    Blue = 'blue',
    Green = 'green',
    Orange = 'orange',
    Purple = 'purple',
}

enum Strategy {
    Canary = 'Canary',
    BlueGreen = 'BlueGreen',
}

const parseImages = (r: RolloutInfo): ImageInfo[] => {
    const images: {[key: string]: ImageInfo} = {};
    const unknownImages: {[key: string]: boolean} = {};

    (r.replicaSets || []).forEach((rs) => {
        (rs.images || []).forEach((img) => {
            const tags: ImageTag[] = [];

            if (rs.canary) {
                tags.push(ImageTag.Canary);
            }
            if (rs.stable) {
                tags.push(ImageTag.Stable);
            }
            if (rs.active) {
                tags.push(ImageTag.Active);
            }
            if (rs.preview) {
                tags.push(ImageTag.Preview);
            }

            if (images[img]) {
                images[img].tags = [...tags, ...images[img].tags];
            } else {
                images[img] = {
                    image: img,
                    tags: tags,
                };
            }

            if (images[img].tags.length === 0) {
                unknownImages[img] = true;
            } else {
                unknownImages[img] = false;
            }
        });
    });

    const imgArray = Object.values(images);
    imgArray.sort((a, b) => {
        return unknownImages[a.image] ? 1 : -1;
    });
    return imgArray;
};

export const Rollout = () => {
    const {name} = useParams<{name: string}>();

    const [rollout, loading] = useWatchRollout(name, true);
    const api = React.useContext(RolloutAPIContext);

    const images = parseImages(rollout);

    for (const img of images) {
        for (const container of rollout.containers) {
            if (img.image === container.image) {
                img.color = ImageColor.Blue;
            }
        }
    }
    const curStep = parseInt(rollout.step, 10) || (rollout.steps || []).length;

    return (
        <div className='rollout'>
            <Helmet>
                <title>{name} / Argo Rollouts</title>
            </Helmet>
            <ThemeDiv className='rollout__toolbar'>
                <ThemeDiv className='rollout__header'>
                    {name} <StatusIcon status={rollout.status as RolloutStatus} />
                </ThemeDiv>
                <React.Suspense fallback={<Spinner />}>
                    <RolloutActions rollout={rollout} />
                </React.Suspense>
            </ThemeDiv>

            <ThemeDiv className='rollout__body'>
                <WaitFor loading={loading}>
                    <ThemeDiv className='info rollout__info'>
                        <div className='info__title'>Summary</div>

                        <InfoItemRow
                            items={{content: rollout.strategy, icon: iconForStrategy(rollout.strategy as Strategy), kind: rollout.strategy?.toLowerCase() as InfoItemKind}}
                            label='Strategy'
                        />
                        <InfoItemRow items={{content: formatTimestamp(rollout.restartedAt), icon: faClock}} label='Last Restarted' />
                        <ThemeDiv className='rollout__info__section'>
                            {rollout.strategy === Strategy.Canary && (
                                <React.Fragment>
                                    <InfoItemRow items={{content: rollout.step, icon: faShoePrints}} label='Step' />
                                    <InfoItemRow items={{content: rollout.setWeight, icon: faBalanceScaleRight}} label='Set Weight' />
                                    <InfoItemRow items={{content: rollout.actualWeight, icon: faBalanceScale}} label='Actual Weight' />{' '}
                                </React.Fragment>
                            )}
                        </ThemeDiv>
                        <ThemeDiv className='rollout__info__section'>
                            <h3>CONTAINERS</h3>
                            {rollout.containers?.map((c) => (
                                <ContainerWidget images={images} key={c.name} container={c} setImage={(image, tag) => api.setRolloutImage(name, c.name, image, tag)} />
                            ))}
                        </ThemeDiv>

                        <h3>IMAGES</h3>
                        <ImageItems images={images} />
                    </ThemeDiv>
                    {rollout.replicaSets && rollout.replicaSets.length > 0 && (
                        <ThemeDiv className='info rollout__info'>
                            <div className='info__title'>Revisions</div>
                            <div style={{marginTop: '1em'}}>
                                {ProcessRevisions(rollout).map((r, i) => (
                                    <RevisionWidget key={i} revision={r} initCollapsed={false} />
                                ))}
                            </div>
                        </ThemeDiv>
                    )}
                    {(rollout.strategy || '').toLowerCase() === 'canary' && rollout.steps && rollout.steps.length > 0 && (
                        <ThemeDiv className='info steps'>
                            <ThemeDiv className='info__title'>Steps</ThemeDiv>
                            <div style={{marginTop: '1em'}}>
                                {rollout.steps.map((step, i) => (
                                    <Step step={step} complete={i < curStep} current={i === curStep} />
                                ))}
                            </div>
                        </ThemeDiv>
                    )}
                </WaitFor>
            </ThemeDiv>
        </div>
    );
};

const iconForStrategy = (s: Strategy) => {
    switch (s) {
        case Strategy.Canary:
            return faDove;
        case Strategy.BlueGreen:
            return faPalette;
    }
};

const ImageItems = (props: {images: ImageInfo[]}) => {
    return (
        <div>
            {props.images.map((img) => {
                let imageItems = img.tags.map((t) => {
                    return {content: t, icon: IconForTag(t)} as InfoItemProps;
                });
                if (imageItems.length === 0) {
                    imageItems = [{icon: IconForTag()}];
                }
                return <InfoItemRow key={img.image} label={<ThemeDiv className={`image image--${img.color || 'unknown'}`}>{img.image}</ThemeDiv>} items={imageItems} />;
            })}
        </div>
    );
};

interface Revision {
    number: number;
    replicaSets: GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1ReplicaSetInfo[];
    experiments: GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1ExperimentInfo[];
    analysisRuns: GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1AnalysisRunInfo[];
}

const ProcessRevisions = (ri: RolloutInfo): Revision[] => {
    if (!ri) {
        return;
    }
    const map: {[key: number]: Revision} = {};

    const emptyRevision = {replicaSets: [], experiments: [], analysisRuns: []} as Revision;

    for (const rs of ri.replicaSets || []) {
        if (!map[rs.revision]) {
            map[rs.revision] = {...emptyRevision};
        }
        map[rs.revision].number = rs.revision;
        map[rs.revision].replicaSets = [...map[rs.revision].replicaSets, rs];
    }

    const revisions: Revision[] = [];
    const prevRn = 0;
    Object.keys(map).forEach((key) => {
        const rn = parseInt(key);
        if (rn > prevRn) {
            revisions.unshift(map[rn]);
        } else {
            revisions.push(map[rn]);
        }
    });

    return revisions;
};

const RevisionWidget = (props: {revision: Revision; initCollapsed?: boolean}) => {
    const {revision, initCollapsed} = props;
    const [collapsed, setCollapsed] = React.useState(initCollapsed);
    const icon = collapsed ? faChevronCircleDown : faChevronCircleUp;
    return (
        <div key={revision.number} style={{marginBottom: '1.5em'}}>
            <ThemeDiv className='revision__header'>
                Revision {revision.number}
                <ThemeDiv className='revision__header__button' onClick={() => setCollapsed(!collapsed)}>
                    <FontAwesomeIcon icon={icon} />
                </ThemeDiv>
            </ThemeDiv>
            {!collapsed && revision.replicaSets.map((rs) => <ReplicaSet key={rs.objectMeta.uid} rs={rs} />)}
        </div>
    );
};

const ContainerWidget = (props: {
    container: GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1ContainerInfo;
    images: ImageInfo[];
    setImage: (image: string, tag: string) => void;
}) => {
    const {container} = props;
    const [editing, setEditing] = React.useState(false);
    const [newImage, setNewImage, newImageInput] = useInput(container.image);

    const switchMode = (editing: boolean) => {
        setEditing(editing);
        setNewImage(container.image);
    };
    return (
        <div style={{margin: '1em 0', display: 'flex', alignItems: 'center'}}>
            {container.name}
            <div style={{marginLeft: 'auto', display: 'flex', alignItems: 'center', height: '2em'}}>
                {!editing ? (
                    <React.Fragment>
                        <InfoItem content={container.image} />
                        <FontAwesomeIcon icon={faPencilAlt} onClick={() => switchMode(true)} style={{cursor: 'pointer', marginLeft: '5px'}} />
                    </React.Fragment>
                ) : (
                    <React.Fragment>
                        <Autocomplete items={props.images.map((img) => img.image)} placeholder='New Image' {...newImageInput} style={{transition: 'width 1s ease'}} />
                        <span style={{marginLeft: '5px'}}>
                            <ActionButton icon={faTimes} action={() => switchMode(false)} />
                        </span>
                        <ActionButton
                            disabled={newImage === '' || newImage.split(':').length < 2}
                            icon={faCheck}
                            action={() => {
                                const split = newImage.split(':');
                                const image = split[0];
                                const tag = split[1];
                                props.setImage(image, tag);
                                setTimeout(() => {
                                    setNewImage('');
                                    setEditing(false);
                                }, 350);
                            }}
                            indicateLoading
                        />
                    </React.Fragment>
                )}
            </div>
        </div>
    );
};

const Step = (props: {step: GithubComArgoprojArgoRolloutsPkgApisRolloutsV1alpha1CanaryStep; complete?: boolean; current?: boolean}) => {
    let icon: IconDefinition;
    let content = '';
    let unit = '';
    if (props.step.setWeight) {
        icon = faWeight;
        content = `Set Weight: ${props.step.setWeight}`;
        unit = '%';
    }
    if (props.step.pause) {
        icon = faPauseCircle;
        if (props.step.pause.duration) {
            content = `Pause: ${props.step.pause.duration}`;
            unit = 's';
        } else {
            content = 'Pause';
        }
    }
    if (props.step.analysis) {
        icon = faChartBar;
    }

    return (
        <div className={`steps__step ${props.complete ? 'steps__step--complete' : ''} ${props.current ? 'steps__step--current' : ''}`}>
            <FontAwesomeIcon icon={icon} /> {content}
            {unit}
        </div>
    );
};

import {Aid, AidList} from "aid-js-sdk"

export const readAidListFromLocalStorage = (): AidList => {
    const defaultUserInfosZip = localStorage.getItem('defaultUserInfosZip');
    const aidsZip = localStorage.getItem('aidsZip');
    if (defaultUserInfosZip === null || aidsZip === null) {
        return new AidList("", "[]");
    }
    return new AidList(defaultUserInfosZip, aidsZip);
}

export const writeAidListToLocalStorage = (aidList: AidList): void => {
    const {
        defaultUserInfosZip,
        aidsZip
    } = aidList.export();

    localStorage.setItem('defaultUserInfosZip', defaultUserInfosZip);
    localStorage.setItem('aidsZip', aidsZip);
}

export const readAid = (aid: string): Aid | null => {
    const aidStr = localStorage.getItem(aid);
    if (aidStr === null) {
        return null;
    }
    return Aid.fromStr(aidStr);
}

export const writeAid = (aid: Aid): void => {
    localStorage.setItem(aid.aid, aid.toStr());
}

export const generateNewAid = (): Aid => {
    const uuid = crypto.randomUUID();
    return new Aid(uuid, new Map(), new Map(), new Map());
}

export const getDefaultAid = (aidList: AidList): Aid | undefined => {
    aidList = readAidListFromLocalStorage();
    if (aidList.aids.length === 0) {
        return undefined
    }
    const targetAid = aidList.aids[0];
    let aid = readAid(targetAid.aid);
    if (aid === null) {
        aid = new Aid(targetAid.aid, new Map(), new Map(), new Map());
    }
    return aid;
}
import {AidList, Aid} from "aid-js-sdk"

export const readAidListFromLocalStorage = (): AidList => {
    const defaultUserInfosZip = localStorage.getItem('defaultUserInfosZip');
    const aidsZip = localStorage.getItem('aidsZip');
    if (defaultUserInfosZip === null || aidsZip === null) {
        return new AidList();
    }
    const aidList = new AidList();
    aidList.defaultUserInfos = new Map<string, string>(JSON.parse(defaultUserInfosZip));
    aidList.aids = JSON.parse(aidsZip);
    return aidList;
}

export const writeAidListToLocalStorage = (aidList: AidList): void => {
    localStorage.setItem('defaultUserInfosZip', JSON.stringify(Array.from(aidList.defaultUserInfos.entries())));
    localStorage.setItem('aidsZip', JSON.stringify(aidList.aids));
}

export const  readAid = (aid: string): Aid | null => {
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
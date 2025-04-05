/**
 * 时间戳 转 date
 * @param nS
 * @returns {string}
 */
function getLocalTime(nS) {
    // 10 需要 * 1000； 13位不做任何操作。
    return new Date(parseInt(nS)).toLocaleString().replace(/:\d{1,2}$/,' ');
}

/* 质朴长存法 by lifesinger */
function pad(num, n) {
    var len = num.toString().length;
    while(len < n) {
        num = "0" + num;
        len++;
    }
    return num;
}

// 通过 date 对象 获取date 时间字符串
function getDateString(d) {
    return d.getFullYear()+'-'+pad(d.getMonth()+1,2)+'-'+pad(d.getDate(), 2)
}

// 通过 date 对象 获取datetime 时间字符串
function getDatetimeString(d) {
    return d.getFullYear()+'-'+pad(d.getMonth()+1,2)+'-'+pad(d.getDate(), 2)+ ' '+ pad(d.getHours())+':'+d.getMinutes()+':'+d.getSeconds()
}

/**
 * 持续时间单位转换
 * @param duration 持续时间
 * @param inFormat 输入单位 us.微秒、ms.毫秒、s.秒、m.分钟、h.小时、d.天
 * @param outFormat 输出格式 latest.表示最近的格式
 * @param lang 语言 en.英文、zh.中文
 * @returns {string}
 */
function durationTransform(duration, inFormat='us', outFormat='latest',lang='en') {
    let langue = {
        us: {zh: "微秒", en: "μs"},
        ms: {zh: "毫秒", en: "ms"},
        s: {zh: "秒", en: "s"},
        m: {zh: "分钟", en: "m"},
        h: {zh: "小时", en: "h"},
        d: {zh: "天", en: "d"},
    };

    if (outFormat == 'latest'){
        if (inFormat == 'us'){ // 因为包最小支持毫秒，所以微秒要转毫秒
            inFormat = 'ms'
            duration /=1000
        }

        let dur = moment.duration(duration, inFormat)
        // 关于单位输出示例
        // moment.duration(6396.528,'ms').seconds()     6
        // moment.duration(6396.528,'ms').asSeconds()   6.396528
        // 特别说明
        // moment.duration(6396.528, 'ms').asMilliseconds()  // 6396.528 输出毫秒单位
        // moment.duration(6396.528, 'ms').milliseconds()    // 396.528 仅输出毫秒的部分

        if (inFormat == 'ms' && duration < 1){
            return duration*1000+langue['us'][lang]
        }else if (dur.asMilliseconds() < 1000){
            return dur.asMilliseconds()+langue['ms'][lang]
        } else if (dur.asSeconds() < 60){
            return dur.asSeconds().toFixed(3)+langue['s'][lang]
        }else if (dur.asMinutes() < 60){
            return dur.asMinutes().toFixed(2) + langue['m'][lang]
        }else if (dur.asHours() < 24){
            return dur.asHours().toFixed(2) + langue['h'][lang]
        }else{
            return (dur.asHours() / 24).toFixed(2) + langue['d'][lang]
        }
    }
    // 跟过功能根据需求实现
    alert("功能未实现...")
}

/**
 * 拖动封装
 * ~~~
 * 注意：使用时需要先引入 SortableJS
 * ~~~
 * @param box 拖动列表包装盒子元素
 * @param list 数据列表；因为json是引用类型，所以内部改变会印象外面。
 * @returns {*} 返回的对象要存储下来，不然会因为执行完销毁而失效。
 * @constructor
 */
function MySortable(box, call){
    return new Sortable(box, {
        handle: ".drag",  // 可拖拽节点
        ghostClass: 'item-drag-current',
        onEnd: event => {
            let nums = box.childNodes.length;
            let newIndex = event.newIndex; // 新位置
            let oldIndex = event.oldIndex; // 原(旧)位置
            // let $label = box.children[newIndex]; // 新节点
            // let $oldLabel = box.children[oldIndex]; // 原(旧)节点
            console.log("拖拽", nums, "old：" + oldIndex, "new：" + newIndex)

            let $label = box.children[newIndex]; // 新节点
            let $oldLabel = box.children[oldIndex]; // 原(旧)节点
            box.removeChild($label);
            if (event.newIndex >= nums) {
                box.insertBefore($label, $oldLabel.nextSibling);
                return;
            }
            if (newIndex < oldIndex) {
                box.insertBefore($label, $oldLabel);
            } else {
                box.insertBefore($label, $oldLabel.nextSibling);
            }

            call(oldIndex, newIndex)

        },
    });
}

/**
 * 一维array数组转树型多维结构
 * @param arrList 数组
 * @param id 主键
 * @param fid 子健
 * @param children 子集名称
 * @returns {*[]|*}
 */
function arrayToTree(arrList, id, fid, children = 'children') {
    let map = []
    arrList.forEach(item => {
        let up = arrList.filter(x => x[id] == item[fid])
        let sub = arrList.filter(x => x[fid] == item[id])
        if (sub.length) item[children] = sub
        if (!up.length) map.push(item)
    })
    return map;
    // if (arrList.length == map.length)return map
    // return arrayToTree(map, id, fid)
}

/**
 * 判断字符串是否为json字符串
 * @param str
 * @returns {boolean}
 */
function isJSON(str) {
    try {
        JSON.parse(str);
        return true;
    } catch (e) {
        return false;
    }
}

/**
 * 解析json字符串
 * @param str
 * @returns {{}|any}
 */
function parseJSON(str, defaultValue = {}){
    try {
        return JSON.parse(str);
    } catch (e) {
        return defaultValue;
    }
}

function stringifyJSON(data){
    return JSON.stringify(data, null, 4)
}

/**
 * json深拷贝
 * @param data
 * @returns {any}
 */
function copyJSON(data) {
    return JSON.parse(JSON.stringify(data))
}

/**
 * 监听输入变化数组追加行
 * ~~~
 * 元素最后一行非空追加行
 * ~~~
 * @param {string|number} e 变化值
 * @param {number} index 变化行号
 * @param {array} arr 原数据
 * @param {object} item 追加的元素，默认{}空对象
 */
function inputChangeArrayPush(e, index, arr, item){
    // console.log(e, index, arr,"-->", arr.length)
    if (e == ""){
        return
    }
    if (item == null){
        item = {}
    }
    if (arr.length-1 <= index){
        arr.push(item) // 追加空元素
    }
    // console.log(e, index, arr, index)
}

/**
 * 删除数组元素
 * ~~~
 * 删除元素行，最后一行不可删除
 * ~~~
 * @param index 删除行号
 * @param arr 原数据
 * @constructor
 */
function arrayDelete(index, arr){
    if ((index+1) >= arr.length){
        return
    }
    arr.splice(index,1)
}

// 跳转回上一页
function backLastPage() {
    if (document.referrer){ // 上一页为当前域，则返回
        window.history.back()
    }else{ // 上一页不是当前域，则跳转到首页
        window.location.href="/index"
    }
}

function getHomePage() {
    window.location.href="/index"
}
// 跳转登录页
function backLoginPage() {
    window.location.href="/login"
}

/**
 * 函数用于解析URL中的hash参数
 * @param {string} url
 * @returns {object}
 */
function getHashParams(url) {
    let hashParams = {};
    // 以'#'为分隔符，取数组的第二个元素（即hash参数）
    url = decodeURIComponent(url)
    let urlParts = url.split('#');
    if (urlParts.length > 1) {
        // 对hash参数字符串进行切割
        let hashs = urlParts[1].split('?')
        let hash = hashs.length > 1 ? hashs[1] : hashs[0]
        let params = hash.split('&');

        params.forEach(function(param) {
            let splitter = param.split('=');
            if (splitter.length > 1) {
                let key = splitter[0]
                let val = splitter[1]
                if (key.length > 2 && key.slice(-2) === '[]'){
                    key = key.slice(0,-2)
                    if (hashParams[key]){
                        hashParams[key].push(val)
                    }else{
                        hashParams[key] = [val]
                    }
                }else{
                    hashParams[key] = val;
                }
            }
        });
    }
    return hashParams;
}

/**
 * 构建查询参数字符串
 * @param {object} data
 */
function buildSearchParamsString(data){
    let p = new URLSearchParams()
    for (let item in data){
        let val = data[item]
        if (Array.isArray(val)){
            if (val.length===0) continue
            val.forEach((item2)=>{
                p.append(item+'[]',item2)
            })
        }else if (val === ''){
            continue
        }else{
            p.append(item, val)
        }
    }
    return p.toString()
}

/**
 * 替换url路径
 * @param {String} hash
 * @param {Object} hashData
 * @param {string} search
 */
function replaceHash(hash, hashData, search=''){
    let str = ''
    if (search !== ''){
        str += '?'+search
    }
    if (hash !== ''){
        str += "#"+hash + "?" + buildSearchParamsString(hashData)
    }else if (hashData){
        str += "#?" + buildSearchParamsString(hashData)
    }

    // console.log("hash",str)
    window.history.replaceState({}, '', str);
}

/**
 * 设置页面名称
 * @param title
 */
function setDocumentTitle(title){
    document.title = title + ' - cronin任务平台'
}

/**
 * 业务枚举
 * @type {{dicEnv: number, dicSqlSource: number, envKey: string}}
 */
const Enum ={
    envKey: "env",
    systemInfoKey: "system_info",
    dicKey: "dic",
    dicEnv: 2,
    dicMsg: 3,
    dicUser: 4,
    dicRole: 5,
    dicTag: 6,
    dicSqlSource: 11,
    dicJenkinsSource: 12,
    dicGitSource: 13,
    dicHostSource: 14,
    dicLogName: 25,   // 日志名称
    dicCmdType: 1001,
    // 1002 git事件
    dicGitEvent: 1002,
    dicSqlDriver: 1003,
    // 任务状态
    dicConfigStatus: 1004,
    // 协议类型
    dicProtocolType: 1005,
    // 1006 接收数据字段
    dicReceiveDataField: 1006,
    // 1007 重试模式
    dicRetryMode: 1007,
    // 1 草稿,停用
    StatusDisable: 1,
    // 5 待审核
    StatusAudited: 5,
    // 6 驳回
    StatusReject: 6,
    // 2 激活
    StatusActive: 2,
    // 3 完成
    StatusFinish: 3,
    // 4 错误
    StatusError: 4,
    // 8 已关闭（预删除）
    StatusClosed: 8,
    // 9 删除
    StatusDelete: 9,
}

function statusTypeName(status){
    switch (Number(status)) {
        case Enum.StatusClosed:
            return 'info';
        case Enum.StatusActive:
        case Enum.StatusAudited:
            return 'primary';
        case Enum.StatusFinish:
            return 'success';
        case Enum.StatusReject:
        case Enum.StatusError:
            return 'warning';
    }
}

/**
 * 获得枚举名称
 * ~~~
 * 用于枚举id转名称
 * ~~~
 * @param enumList
 * @param id
 */
function getEnumName(enumList, id){
    let data = enumList.filter((item)=>{return id == item.id})
    return data.length ? data[0].name : ''
}

/**
 * 网络请求
 * @type {{post: api.post, get: api.get}}
 */
var api = {
    baseUrl : window.location.protocol+"//"+window.location.host,
    env : {},
    token: '',

    _getToken(){
        if (!this.token){
            this.token = cache.getToken()
        }
        return this.token
    },
    /**
     * 内网get请求
     * @param path string 请求路径
     * @param param object 请求参数 将拼url参数
     * @param success func 响应结果
     * @param setting object 请求设置
     */
    innerGet: function (path,param, success, setting) {
        let url =  this.baseUrl+path
        let async = true
        if (setting && typeof setting.async == "boolean"){
            async = setting.async
        }

        let header = {
            'env': this.getEnv().env,
            'Authorization': this._getToken()
        }
        this.ajax('get', url, param, header,success, async)
    },

    /**
     * 内网post请求
     */
    innerPost: function (path, param, success) {
        let url =  this.baseUrl+path
        param = JSON.stringify(param)
        let header = {
            'Content-Type': 'application/json',
            'env': this.getEnv().env,
            'Authorization': this._getToken()
        }
        this.ajax('post', url, param, header, success)
    },

    innerGetFile: function(path){
        let features = "height=240, width=400, top=50, left=50, toolbar=no, menubar=no,scrollbars=no,resizable=no, location=no, status=no,toolbar=no";
        window.open(this.baseUrl+path, "_blank", features)
    },

    /**
     * 枚举列表
     * @param types array 枚举key列表
     * @param callback function() 回调函数
     * @param reload bool 是否强制重载
     * @returns {Promise<void>}
     */
    dicList: function (types, callback, reload= false) {

        let list = JSON.parse(localStorage.getItem(Enum.dicKey)) ?? {}
        let reply = {}
        let queryInfo = []
        // 从存储查看要请求的枚举是否存在 存在从存储取出
        types.forEach((element, index) => {
            if (list[element] && !reload) {
                reply[element] = list[element]
            } else {
                queryInfo.push(element)
            }
        })

        // 不存在存储的枚举继续查询
        if (queryInfo.length > 0) {
            // 查询不存在缓存的部分调取一次请求，并缓存；如果都存在缓存中就无需再次发起请求。
            this.innerGet("/foundation/dic_gets", {"types":queryInfo.join(',')}, res =>{
                if (res.status) {
                    queryInfo.forEach(element => {
                        if (res.data.maps[element]){
                            reply[element] = res.data.maps[element].list
                            list[element] = reply[element]
                        }
                    })
                    window.localStorage.setItem(Enum.dicKey, JSON.stringify(list))
                } else {
                    console.log("枚举错误",res)
                    alert('枚举错误:' + res.message)
                }
            }, {async:false})

        }
        callback(reply)
    },
    /**
     * 移除dic缓存
     * @param types array 枚举key列表
     */
    dicDel: function (types) {
        let list = JSON.parse(localStorage.getItem(Enum.dicKey)) ?? {}
        let reply = {}
        // 从存储查看要请求的枚举是否存在 存在从存储取出
        types.forEach((type, index) => {
            if (list[type]) {
                delete list[type]
            }else if (type == -1){
                list = []
            }
        })

        window.localStorage.setItem(Enum.dicKey, JSON.stringify(list))
    },

    // 系统信息
    // 信息具有缓存功能
    systemInfo: function (callback, reload=false) {
        // 不重载，且缓存存在，使用缓存数据
        if (!reload){
            let list = JSON.parse(localStorage.getItem(Enum.systemInfoKey)) ?? {}
            if (list){
                return callback(list)
            }
        }

        this.innerGet("/foundation/system_info", null,(res) =>{
            if (!res.status){
                return callback(res);
            }

            let envStr = localStorage.getItem(Enum.envKey)
            if (envStr != null && envStr != ""){
                let env = JSON.parse(envStr)
                if (env.env != "" && env.env_name != ""){
                    res.data.env = env.env
                    res.data.env_name = env.env_name
                }
            }else {
                this.setEnv(res.data.env, res.data.env_name)
            }
            localStorage.setItem(Enum.systemInfoKey, JSON.stringify(res.data))
            callback(res.data)
        }, {async:false})
    },

    setEnv(key, name){
        this.env = {env:key, env_name:name}
        localStorage.setItem(Enum.envKey, JSON.stringify(this.env))
    },

    getEnv(){
        if (!this.env.env){
            let envStr = localStorage.getItem(Enum.envKey)
            let env = JSON.parse(envStr)
            if (env == null || env.env == ""){
                this.systemInfo((res)=>{
                    if (!res.status){
                        alert(res.message);
                    }
                })
            }else{
                this.env = env
            }
        }
        return this.env
    },

    ajax: function (method, url, data, header, callback, async=true, dataType='json') {
        $.ajax({
            'url': url,
            'data': data,
            'type': method,
            'headers': header,
            'dataType': dataType,
            'async': async,
            'success': res => {
                if (res.code == '000000'){
                    res.status = true
                }else if (res.code == '999909'){
                    // 重定向到登录
                    return backLoginPage()
                }else{
                    res.status = false
                }
                return callback(res)
            },
            'error': res =>{
                let temp = {}
                if (res.status == 0){
                    temp.code = res.status
                    temp.message = res.statusText
                }else if(Object.keys(res.responseJSON).length > 0){
                    temp = res.responseJSON
                }else {
                    temp.code = res.status
                    temp.message = res.status + " " + res.statusText
                }
                temp.status = false
                return callback(temp)
            }
        });
    }
}

// 缓存信息
var cache = {

    /**
     * 获得环境
     * @returns {any|{}}
     */
    getEnv(){
        let list = JSON.parse(localStorage.getItem('env')) ?? {}
        return list
    },

    /**
     * 设置环境
     * @param val
     */
    setEnv(val){
        let str = ''
        if (typeof val === 'string'){
            str = val
        }else if (typeof val === 'object'){
            str = JSON.stringify(val)
        }

        localStorage.setItem('user', str)
    },

    /**
     * 获得令牌
     * @returns {string}
     */
    getToken(){
        let str = localStorage.getItem('token')
        return str
    },
    /**
     * 设置令牌
     * @param str
     */
    setToken(str){
        localStorage.setItem('token', str)
    },
    /**
     * 获取用户信息
     * @returns {any|{}}
     */
    getUser(){
        let list = JSON.parse(localStorage.getItem('user')) ?? {}
        return list
    },
    /**
     * 设置用户信息
     * @param user
     */
    setUser(user){
        let str = ''
        if (typeof user === 'string'){
            str = user
        }else if (typeof user === 'object'){
            str = JSON.stringify(user)
        }

        localStorage.setItem('user', str)
    },
    /**
     * 获得权限节点信息
     * @returns {any|{}}
     */
    getMenu(){
        let list = JSON.parse(localStorage.getItem('menus')) ?? {}
        return list
    },
    /**
     * 设置权限节点信息
     * @param menus
     */
    setMenu(menus){
        let str = ''
        if (typeof menus === 'string'){
            str = menus
        }else if (typeof menus === 'object'){
            str = JSON.stringify(menus)
        }

        localStorage.setItem('menus', str)
    },
    /**
     * 设置权限标记
     * @param tags
     */
    setAuthTags(tags){
        let str = ''
        if (typeof tags === 'string'){
            str = tags
        }else if (typeof tags === 'object'){
            str = JSON.stringify(tags)
        }

        localStorage.setItem('auth_tags', str)
    },
    /**
     * 获得权限标记
     * @returns {any|{}}
     */
    getAuthTags(){
        let list = JSON.parse(localStorage.getItem('auth_tags')) ?? {}
        return list
    },
    empty(){
        localStorage.removeItem('dic')
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        localStorage.removeItem('system_info')
        localStorage.removeItem('menus')
    },
    delDic(){
        localStorage.removeItem('dic')
    }
}
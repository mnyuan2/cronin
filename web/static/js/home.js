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
        } else if (dur.seconds() < 60){
            return Math.round(dur.asMilliseconds())/1000+langue['s'][lang]
        }else if (dur.minutes() < 60){
            return dur.seconds() / 60 + langue['m'][lang]
        }else if (dur.hours() < 24){
            return dur.minutes() / 60 + langue['h'][lang]
        }else{
            return dur.hours() / 24 + langue['d'][lang]
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
        let sit = arrList.filter(x => x[fid] == item[id])
        if (sit.length) item[children] = sit
        if (!(up.length && !sit.length)) map.push(item)
    })
    return map.length > 0 ? [map[0]] : map;
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
    dicSqlSource: 11,
    dicJenkinsSource: 12,
}

/**
 * 网络请求
 * @type {{post: api.post, get: api.get}}
 */
var api = {
    baseUrl : window.location.protocol+"//"+window.location.host,
    env : {},

    /**
     * 内网get请求
     * @param path string 请求路径
     * @param param object 请求参数 将拼url参数
     * @param success func 响应结果
     */
    innerGet: function (path,param, success) {
        let url =  this.baseUrl+path
        // if (param){
        //     url += "?"
        //     for (i in param){
        //         if (param[i]){
        //             url += i+"="+param[i]+"&"
        //         }
        //     }
        //     url = url.slice(0,-1)
        // }

        let header = {
            'env': this.getEnv().env
        }
        this.ajax('get', url, param, header,success)
    },

    /**
     * 内网post请求
     */
    innerPost: function (path, param, success) {
        let url =  this.baseUrl+path
        param = JSON.stringify(param)
        let header = {
            'Content-Type': 'application/json',
            'env': this.getEnv().env
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
    dicList: async function (types, callback, reload= false) {

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
            // 这里的请求必须等待结果
            await new Promise((resolve, reject) => {
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
                    resolve()
                })
            })
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

        this.ajax('get', this.baseUrl+"/foundation/system_info", null,null, (res) =>{
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
        }, false)
    },

    setEnv(key, name){
        this.env = {env:key, env_name:name}
        localStorage.setItem(Enum.envKey, JSON.stringify(this.env))
    },

    getEnv(){
        if (this.env.length == 0 || !this.env.env){
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
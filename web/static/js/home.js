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
 * 业务枚举
 * @type {{dicEnv: number, dicSqlSource: number, envKey: string}}
 */
const Enum ={
    envKey: "env",
    systemInfoKey: "system_info",
    dicKey: "dic",
    dicSqlSource: 1,
    dicEnv: 2,
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
            if (list[element]) {
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
                if(Object.keys(res.responseJSON).length > 0){
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
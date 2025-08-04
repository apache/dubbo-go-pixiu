/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/**弹框底部按钮点击事件命令名称 */
export const cmds = {
    /**返回 */
    goBack: 'goBack',
    /**上一页 */
    prev: 'prev',
    /**下一页 */
    next: 'next',
    /**新增用户 */
    saveCustInfo: 'saveCustInfo',
    /**上传用户档案 */
    uplCustArch: 'uplCustArch',
    /**车辆新增 */
    saveVehicleInfo: 'saveVehicleInfo',
    /**车辆档案上传 */
    uplVehicleArch: 'uplVehicleArch',
    /**读卡 */
    readCard: 'readCard',
    /**卡发行 */
    cardIssue: 'cardIssue',
    /**读标签 */
    readObu: 'readObu',
    /**标签发行 */
    obuIssue: 'obuIssue',
    /**历史受理单 */
    historyReceipts: 'historyReceipts',
    /**生成受理单 */
    createReceipts: 'createReceipts',
    /**上传受理单 */
    uploadReceipts: 'uploadReceipts',
    /**打印受理单 */
    prtReceipts: 'prtReceipts',
    /**标签注销 */
    labelCancellation: 'labelCancellation',
    /**标签打印 */
    labelPrint: 'labelPrint',
    /**车辆注销 */
    vehicleCancellation: 'vehicleCancellation',
    /**强制注销 */
    compulsoryCancellation: 'compulsoryCancellation',
    /**撤销签约*/
    unSign: 'unSign',
    /**卡片检测*/
    checkCpu: 'checkCpu',
    /**标签检测*/
    checkObu: 'checkObu',
    /***/
    webIssue: 'webIssue',
    /**储值资金*/
    storedValue:'storedValue',
    mRegistration:'mRegistration',
    /**特殊签约*/
    specialSign: 'specialSign',
    /**卡片冻结*/
    frozen: 'frozen',
    /**卡片解冻*/
    unFrozen: 'unFrozen',
    /**打印冻结解冻业务单*/
    prtFrozen: 'prtFrozen',
    /**卡片挂失*/
    loss: 'loss',
    /**卡片解挂*/
    unLoss: 'unLoss',
    /**打印挂失解挂业务单*/
    prtLoss: 'prtLoss',
    /**补领换卡*/
    reCard: 'reCard',
    /**卡片解锁*/
    unlock: 'unlock',
    /**pos充值*/
    posRecharge:'posRecharge',
    /**账户明细*/
    accountInfo:'accountInfo',

    /**卡片注销*/
    cardCancellation:'cardCancellation',
    /**退款单打印*/
    refundPrinting:'refundPrinting',
    /**无卡销卡*/
    noCard:'noCard',
    /**有卡销卡*/
    hasCard:'hasCard',
    /**圈存异常人工处理*/
    manualHandle:'manualHandle',
    /**新卡校验*/
    newCardCheck:'newCardCheck',
    /**原卡销卡*/
    originalCard:'originalCard',
    /**业务单打印*/
    BusinessOrderPrinting:'BusinessOrderPrinting',
    /**补领换卡 -》发卡*/
    newCardissue:'newCardissue',

}
/**弹框底部按钮 */
export const ftBtns = {
    //用户信息
    customerInfo: [
        {
            text: "返回",
            cls: "gray-btn space30 arrowleft",
            cmd: cmds.goBack
        },
        {
            text: "保存",
            cls: "blue-btn space30 flex2",
            cmd: cmds.saveCustInfo
        },
        {
            text: "用户档案",
            cls: "green-btn  arrowright",
            cmd: cmds.next
        }
    ],
    //用户档案
    customerArch: [
        {
            text: "用户信息",
            cls: "gray-btn  arrowleft",
            cmd: cmds.prev
        },
        {
            text: "保存",
            cls: "blue-btn space30 flex2",
            cmd: cmds.uplCustArch
        },
        {
            text: "车辆信息",
            cls: "green-btn   arrowright",
            cmd: cmds.next
        }
    ],
    //车辆信息
    vehicleInfo: [
        {
            text: "用户档案",
            cls: "gray-btn  arrowleft",
            cmd: cmds.prev
        },
        {
            text: "保存",
            cls: "blue-btn space30 flex2",
            cmd: cmds.saveVehicleInfo
        },
        {
            text: "车辆档案",
            cls: "green-btn  arrowright",
            cmd: cmds.next
        }
    ],
    //车辆档案
    vehicleArch: [
        {
            text: "车辆信息",
            cls: "gray-btn  arrowleft",
            cmd: cmds.prev
        },
        {
            text: "上传",
            cls: "blue-btn space30 flex2",
            cmd: cmds.uplVehicleArch
        },
        {
            text: "卡片发行",
            cls: "green-btn  arrowright",
            cmd: cmds.next
        }
    ],
    //卡片发行
    cardIssue: [
        {
            text: "车辆档案",
            cls: "gray-btn  arrowleft",
            cmd: cmds.prev
        },
        {
            text: "读卡",
            cls: "green-btn space30",
            cmd: cmds.readCard
        },
        {
            text: "卡发行",
            cls: "blue-btn",
            cmd: cmds.cardIssue
        },
        {
            text: "标签发行",
            cls: "green-btn  arrowright",
            cmd: cmds.next
        }
    ],
    //标签发行
    obuIssue: [
        {
            text: "卡片发行",
            cls: "gray-btn  arrowleft",
            cmd: cmds.prev
        },
        {
            text: "读 标 签",
            cls: "green-btn",
            cmd: cmds.readObu
        },
        {
            text: "标签发行",
            cls: "blue-btn",
            cmd: cmds.obuIssue
        },
        {
            text: "受 理 单",
            cls: "green-btn  arrowright",
            cmd: cmds.next
        }
    ],
    //受理单
    receipts: [
        {
            text: "历史受理单",
            cls: "gray-btn ",
            cmd: cmds.historyReceipts
        },
        {
            text: "生成受理单",
            cls: "green-btn",
            cmd: cmds.createReceipts
        },
        {
            text: "上传受理单",
            cls: "blue-btn",
            cmd: cmds.uploadReceipts
        },
        {
            text: "打印",
            cls: "green-btn ",
            cmd: cmds.prtReceipts
        }
    ],
    //卡片检测
    checkCard: [
        {
            text: "读卡",
            cls: "green-btn space30  ",
            cmd: cmds.readCard
        }
    ],
    //车辆注销
    vehicleCancellation: [
        {
            text: "车辆注销",
            cls: "green-btn space30  ",
            cmd: cmds.vehicleCancellation
        }
    ],
    //车辆注销
    compulsoryCancellation: [
        {
            text: "车辆强制注销",
            cls: "green-btn space30  ",
            cmd: cmds.compulsoryCancellation
        }
    ],
    //标签检测
    checkObu: [
        {
            text: "读标签",
            cls: "green-btn space30  ",
            cmd: cmds.readObu
        },
    ],
    unSign: [
        {
            text: "撤销签约",
            cls: "green-btn space30 radius-r-b-6 radius-l-b-6",
            cmd: cmds.unSign
        }
    ],
    specialSign: [
        {
            text: "保存",
            cls: "green-btn space30 radius-r-b-6 radius-l-b-6",
            cmd: cmds.specialSign
        }
    ],
    labelCancellation: [{
            text: "打印",
            cls: "green-btn space30 radius-l-b-6",
            cmd: cmds.labelPrint
        }, {
            text: "标签注销",
            cls: "green-btn space30 radius-r-b-6",
            cmd: cmds.labelCancellation
    }],
    posRecharge:[{
            text: "撤销",
            cls: "green-border-btn space30 radius-l-b-6 is-disabled",
            cmd: 'revoke'
        }, {
            text: "充值",
            cls: "green-btn space30 radius-r-b-6",
            cmd: 'recharge'
    }],
    posCorrect:[{
        text: "确认冲正",
        cls: "green-btn space30 radius-r-b-6",
        cmd: 'correct'
    }],
    //卡片冻结解冻
    cardFrozen: [
        {
            text: "解除冻结",
            cls: "gray-btn",
            cmd: cmds.unFrozen
        },
        {
            text: "打印",
            cls: "blue-btn space30",
            cmd: cmds.prtFrozen
        },
        {
            text: "卡片冻结",
            cls: "green-btn",
            cmd: cmds.frozen
        }
    ],
    //补领换卡
    reCard: [

    ],

    //卡片挂失解挂
    cardLoss: [
        {
            text: "解除挂失",
            cls: "gray-btn",
            cmd: cmds.unLoss
        },
        {
            text: "打印",
            cls: "blue-btn space30",
            cmd: cmds.prtLoss
        },
        {
            text: "卡片挂失",
            cls: "green-btn",
            cmd: cmds.loss
        }
    ],
    cardUnlock: [
        {
            text: "卡片解锁",
            cls: "green-btn space30 radius-r-b-6 radius-l-b-6",
            cmd: cmds.unlock
        }
    ],
    //卡片注销
    cardCancellation: [
        {
            text: "退款单打印",
            cls: "gray-btn flex2",
            cmd: cmds.refundPrinting
        },
        {
            text: "无卡销卡",
            cls: "blue-btn  flex2",
            cmd: cmds.noCard
        },
        {
            text: "读卡",
            cls: "green-btn  flex2",
            cmd: cmds.readCard
        },
        {
            text: "有卡销卡",
            cls: "green-btn flex2",
            cmd: cmds.hasCard
        }
    ],
    //卡片注销2
    cardCancellationTwo: [
        {
            text: "打印",
            cls: "gray-btn flex2",
            cmd: cmds.refundPrinting
        },
        {
            text: "无卡销卡",
            cls: "blue-btn  flex2",
            cmd: cmds.noCard
        },
        {
            text: "读卡",
            cls: "green-btn  flex2",
            cmd: cmds.readCard
        },
        {
            text: "有卡销卡",
            cls: "green-btn flex2",
            cmd: cmds.hasCard
        }
    ],
    manualHandle:[{
        text: "人工异常处理",
        cls: "green-btn space30 radius-r-b-6",
        cmd: cmds.manualHandle
    }],
    newCardCheck:[{
        text: "读取新卡",
        cls: "green-btn space30",
        cmd: cmds.readCard
    }, {
        text: "新卡校验",
        cls: "blue-btn space30",
        cmd: cmds.newCardCheck
    },],
    businessOrderPrinting:[{
        text: "新卡校验",
        cls: "gray-btn space30 arrowleft",
        cmd: cmds.newCardCheck
    }, {
        text: "业务单打印",
        cls: "blue-btn space30",
        cmd: cmds.BusinessOrderPrinting
    }, {
        text: "原卡销卡",
        cls: "blue-btn space30 arrowright",
        cmd: cmds.originalCard
    }],
    deleteCard:[{
        text: "业务单打印",
        cls: "gray-btn  arrowleft",
        cmd: cmds.BusinessOrderPrinting
    }, {
        text: "无卡销卡",
        cls: "green-btn",
        cmd: cmds.noCard
    },{
        text: "读取原卡",
        cls: "blue-btn",
        cmd: cmds.readCard
    }, {
        text: "原卡销卡",
        cls: "green-btn",
        cmd: cmds.hasCard
    }],

    //补领换卡  -  新卡发行
    newCardissue: [
        {
            text: "读取新卡",
            cls: "green-btn space30",
            cmd: cmds.readCard
        },
        {
            text: "新卡发行",
            cls: "blue-btn",
            cmd: cmds.cardIssue
        }
    ],
}

export default ftBtns
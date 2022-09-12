import http from "k6/http";
import { group, check } from "k6";
import { scenario } from 'k6/execution';
const body = "{\"destination\":\"Ufab9da979accde2927154cfd153f2fe9\",\"events\":[{\"type\":\"message\",\"message\":{\"type\":\"text\",\"id\":\"15955241648900\",\"text\":\"おみくじ\"},\"webhookEventId\":\"01G15S4GEN3PYN1QAVKR706P8K\",\"deliveryContext\":{\"isRedelivery\":false},\"timestamp\":1650535317544,\"source\":{\"type\":\"user\",\"userId\":\"MOCK_USER_DAS08\"},\"replyToken\":\"f5c0089d22cc425db33d2c7c86f94dff\",\"mode\":\"active\"}]}";

export const options = {
    scenarios: {
        'scenarios': {
            executor: 'shared-iterations',
            vus: 2,
            iterations: 2,
            maxDuration: '30s',
        },
    },
};

export function setup() {
    console.log("setup");
}

export default function () {
    const authParams = {
        headers: {
            'Content-Type': 'application/json',
            'X-Line-Signature': 'Xbuzdtme8zyXfjM0Jpc/WsKgxwqGI7Yc1GfrQjFLUhN=',
        }
    };
    group("/xxxxx/xxxxx", () => {
        const response = http.post(`${__ENV.RAKUTAN_BOT_API}/callback`, body, authParams);
        check(response, {
            'is status 200': (r) => r.status === 200,
        });
    });
}

export function teardown() {
    console.log("teardown");
}
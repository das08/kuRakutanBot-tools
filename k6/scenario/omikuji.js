import http from "k6/http";
import { group, check } from "k6";
import { scenario } from 'k6/execution';
const body = "{\"destination\":\"Ufab9da979accde2927154cfd153f2fe9\",\"events\":[{\"type\":\"message\",\"message\":{\"type\":\"text\",\"id\":\"15955241648900\",\"text\":\"おみくじ\"},\"webhookEventId\":\"01G15S4GEN3PYN1QAVKR706P8K\",\"deliveryContext\":{\"isRedelivery\":false},\"timestamp\":1650535317544,\"source\":{\"type\":\"user\",\"userId\":\"MOCK_USER_DAS08\"},\"replyToken\":\"f5c0089d22cc425db33d2c7c86f94dff\",\"mode\":\"active\"}]}";

export const options = {
    // シナリオ設定
    scenarios: {
        'scenarios': {
            executor: 'shared-iterations', // 複数のVUでiterationを共有
            vus: 2, // 同時接続数
            iterations: 2, // シナリオの総反復回数
            maxDuration: '30s', // 試験の実行時間
        },
    },
    // テスト対象システムの期待性能の合格・不合格を判断する閾値
    // thresholds: {
    //     http_req_failed: ['rate<0.01'], // エラーが1%を超えない
    //     http_req_duration: ['p(95)<500'], // リクエストの95%は200ms以下であること
    // },
};

export function setup() {
    // テスト実行前に実行
    console.log("setup");
}

export default function () {

    // アクセストークンの必要な通信実行
    const authParams = {
        headers: {
            'Content-Type': 'application/json',
            'X-Line-Signature': 'Xbuzdtme8zyXfjM0Jpc/WsKgxwqGI7Yc1GfrQjFLUhN=',
        }
    };
    group("/xxxxx/xxxxx", () => {
        const response = http.post('http://host.docker.internal:8081/callback', body, authParams);
        check(response, {
            'is status 200': (r) => r.status === 200,
        });
    });
}

export function teardown() {
    // テスト実行後に実行
    console.log("teardown");
}
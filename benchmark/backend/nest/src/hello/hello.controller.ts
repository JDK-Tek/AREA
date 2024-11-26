import { Controller, Get, Param, Query } from '@nestjs/common';

function calculate(t0: number[], t1: number[], it: number): [boolean, number, number] {
    const start = performance.now();
    let vel = { x: 0, y: 0, z: 0 };
    let x = 0;
    let k = 0;
    let angle = 0;
    let res = true;

    for (let n = 0; n < it; ++n) {
        x += n;
        if ((t1[2] == t0[2]) || (t0[2] > 0 && t1[2] < 0) || (t0[2] < 0 && t1[2] > 0)) {
            res = false;
            continue;
        }
        vel.x = t1[0] - t0[0];
        vel.y = t1[1] - t0[1];
        vel.z = t1[2] - t0[2];
        k = Math.sqrt(vel.x ** 2 + vel.y ** 2 + vel.z ** 2);
        angle = Math.asin(vel.z / k);
        angle = Math.floor(angle * -180 / Math.PI * 100) / 100;
    }
    console.log(x);
    return [res, angle, performance.now() - start];
}

@Controller('hello')
export class HelloController {
    @Get(':name')
    getHello(
        @Param('name') name: string,
        @Query('t0') t0str: string,
        @Query('t1') t1str: string
    ): string {
        const t0 = JSON.parse(t0str);
        const t1 = JSON.parse(t1str);
        const [success, res, ms] = calculate(t0, t1, 1_000_000);
        if (!success)
            return `Hello ${name}! your ball wont reach, computed in ${ms.toFixed(2)}ms`
        return `Hello ${name}! your incidence angle is ${res.toFixed(2)}, computed in ${ms.toFixed(2)}ms`
    }
}

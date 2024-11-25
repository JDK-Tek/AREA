import { Controller, Post, Body, BadRequestException } from '@nestjs/common';
import { createHash } from 'crypto';
import * as jwt from 'jsonwebtoken';

// soon in env variable
const expirationTime = Math.floor(Date.now() / 1000) + 30 * 60;

const checkForCredentials = (email: string, password: string): string => {
    const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailRegex.test(email))
        throw new BadRequestException('Invalid email format');
    
    if (!/[a-z]/.test(password) || !/[A-Z]/.test(password)
        || !/[0-9]/.test(password) || !/[@!\\/?:+-_$<>#]/.test(password)
        || password.length < 8)
        throw new BadRequestException("Password need to have 8 characters with\
            at least number, uppercase, lowercase and special characters")
    
    const hashedPassword = createHash('sha256').update(password).digest('hex');
    return hashedPassword;
}

@Controller('api')
export class ApiController {
    @Post('register')
    register(@Body() body: { email: string; password: string }) {
        const { email, password } = body;
        
        if (!email || !password)
            throw new BadRequestException('An email and password are needed');
        
        const hashedPassword = checkForCredentials(email, password);

        console.log('Email:', email);
        console.log('Password:', hashedPassword);

        // todo: store it in db

        const payload = {
            email: email,
            exp: expirationTime,
        };

        const token = jwt.sign(payload, hashedPassword);

        return {  token: token };
    }

    @Post('login')
    login(@Body() body: { email: string; password: string }) {
        const { email, password } = body;
        
        if (!email || !password)
            throw new BadRequestException('An email and password are needed');
        
        const hashedPassword = checkForCredentials(email, password);

        console.log('Email:', email);
        console.log('Password:', hashedPassword);

        // todo: check if password & email are good

        const payload = {
            email: email,
            exp: expirationTime,
        };

        const token = jwt.sign(payload, hashedPassword);

        return {  token: token };
    }
}

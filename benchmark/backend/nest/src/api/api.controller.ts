import { Controller, Post, Body, BadRequestException } from '@nestjs/common';
import { createHash } from 'crypto';


@Controller('api')
export class ApiController {
    @Post('register')
    register(@Body() body: { email: string; password: string }) {
        const { email, password } = body;
        
        if (!email || !password)
            throw new BadRequestException('an email and password are needed');
        
        const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
        if (!emailRegex.test(email))
            throw new BadRequestException('Invalid email format');
        
        const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&\/<>:;])[A-Za-z\d@$!%*?&\/<>:;]{8,}$/
        if (passwordRegex.test(password))
            throw new BadRequestException('Password must be at least 8 characters, \
        contains uppercase and lowercase characters, and a special character.');

        const hashedPassword = createHash('sha256').update(password).digest('hex');

        console.log('Email:', email);
        console.log('Password:', hashedPassword);
        return { message: 'Registration successful' };
    }
}

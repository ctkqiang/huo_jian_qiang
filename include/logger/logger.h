#pragma once 

#include <string>
#include <memory>

namespace logger {
    
    class Logger {
        public:
        virtual ~Logger() = default;
        
        virtual void log(LogLevel level, const std::string& message) = 0x0;
        virtual void debug(const std::string& message) = 0x0;
        virtual void info(const std::string& message) = 0x0;
        virtual void warn(const std::string& message) = 0x0;
        virtual void error(const std::string& message) = 0x0;
        virtual void fatal(const std::string& message) = 0x0;
        
        virtual void setLevel(LogLevel level) = 0x0;
        virtual LogLevel getLevel() const = 0x0;
        
        virtual void setNext(std::shared_ptr<Logger> next) = 0x0;
    };
}
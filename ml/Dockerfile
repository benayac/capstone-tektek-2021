FROM python:3.9

# 
WORKDIR /app

# 
COPY /ml/requirements.txt /app/requirements.txt

# 
RUN pip install -r /app/requirements.txt

# 
COPY /ml ./

# 
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8001"]
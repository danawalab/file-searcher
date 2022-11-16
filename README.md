# File-Searcher
대용량 파일을 색인/검색할 수 있는 검색 도구입니다. indexing & search large text files.

## Specification

	
1. 텍스트 파일에 대해 색인 및 검색 기능을 제공한다. ( Provides index and search capabilities for text files )
2. 수집되는 파일에 대해 실시간으로 색인한다. ( Indexing files collected in real time )
3. 서비스는 설정을 정의할 수 있고, 문서 포맷별 파싱 방법을 설정할 수 있다. ( The service can define settings and set parsing methods for each document format )
4. 날짜별, 구분별 조회 조건을 이용한 검색이 가능하다. ( The service can define settings and set parsing methods for each document formatIt is possible to search using inquiry conditions by date and classification )
5. 색인 대상은 타겟 디렉토리에 존재하는 파일들이다. ( The indexing targets are files that exist in the target directory )
6. 메모리에 적재된 색인 파일을 주기적으로 삭제 가능하다. ( Index files loaded in memory can be periodically deleted )
7. API를 통해 다른 클라이언트에서 호출해서 결과를 받을 수 있다. ( You can call from other clients through API and receive results )
8. 로그를 통해 실시간으로 메모리에 적재한 문서와 메모리 사용량을 확인할 수 있다. ( You can check the documents loaded in memory and the memory usage in real time through the log )

## Settings

설정 파일(/config/config.yaml)을 시스템에 맞게 설정합니다.

- 설정 파일 셋팅(server)
1. port : 프로그램이 동작할 포트를 입력합니다.
2. epfilepath: 텍스트 파일이 위치한 경로를 입력합니다.
3. tempfilepath : 임시 청크파일이 위치할 경로를 입력합니다.
4. indexInterval : 색인 간격을 입력합니다.
5. fileDeletePeriodDay : 파일을 삭제할 날짜 간격을 입력합니다.
6. cpuCore : 시스템에서 사용할 CPU 코어 수를 셋팅합니다.
7. workers : 동시에 실행할 경량 쓰레드 갯수를 제어합니다.

- 설정 파일 셋팅(shop)
1. shop: 해당하는 업체의 파싱 방법을 설정합니다.
2. delimiter: 색인할 대상의 구분자를 설정합니다.
3. idPostion: 색인할 대상의 필드 위치를 입력합니다.
4. header: 해당 문서의 헤더가 존재하는지 구분합니다.
5. custom: 싱글라인, 멀티라인 문서인지 구분합니다.

## Pattern example

1. 파일 명 (file name)
```
EE311-1667786400-2.txt
(Company Name)-(Timestamp)-(Overall/Partial part status).txt
```

2. 싱글 라인 파일 (single line)
```
1000520921557	송월타올 오가닉 고중량 세면타월180g 3P 그레이	16900	16900	16900	http://emart.ssg.com/item/itemView.ssg?itemId=1000520921557&siteNo=6001&ckwhere=danawa&appPopYn=n&utm_medium=PCS&utm_source=danawa&utm_campaign=danawa_pcs	http://m.emart.ssg.com/item/itemView.ssg?itemId=1000520921557&siteNo=6001&ckwhere=danawa1&appPopYn=n&utm_medium=PCS&utm_source=danawa&utm_campaign=danawa_pcs	https://sep-item.ssgcdn.com/57/15/92/item/1000520921557_i1_350.jpg		가구/인테리어	침구/커튼	타월	수건/목욕타월			신상품	N	N	N		N	S	8801346015404			송월타월				지금 멤버십 가입하면 1개월 무료!		N									1	0	3000								D	2022-11-07 10:10:48
1000010999599	양치컵SET_조리개포장	4700	4700	4700	http://emart.ssg.com/item/itemView.ssg?itemId=1000010999599&siteNo=6001&ckwhere=danawa&appPopYn=n&utm_medium=PCS&utm_source=danawa&utm_campaign=danawa_pcs	http://m.emart.ssg.com/item/itemView.ssg?itemId=1000010999599&siteNo=6001&ckwhere=danawa1&appPopYn=n&utm_medium=PCS&utm_source=danawa&utm_campaign=danawa_pcs	https://sep-item.ssgcdn.com/99/95/99/item/1000010999599_i1_350.jpg		생활/자동차/공구/성인	욕실용품 	구강용품	구강관리용품			신상품	N	N	N		N				WD730	위덴(wedent)	위덴 치아랑(주)	한국		지금 멤버십 가입하면 1개월 무료!		N									1	124	3000								U	2022-11-07 10:10:48
1000011350898	VIOlight 개인용/휴대용 UV 칫솔 살균기	22000	22000	22000	http://emart.ssg.com/item/itemView.ssg?itemId=1000011350898&siteNo=6001&ckwhere=danawa&appPopYn=n&utm_medium=PCS&utm_source=danawa&utm_campaign=danawa_pcs	http://m.emart.ssg.com/item/itemView.ssg?itemId=1000011350898&siteNo=6001&ckwhere=danawa1&appPopYn=n&utm_medium=PCS&utm_source=danawa&utm_campaign=danawa_pcs	https://sep-item.ssgcdn.com/98/08/35/item/1000011350898_i1_350.jpg		생활/자동차/공구/성인	욕실용품 	구강용품	구강관리용품			신상품	N	N	N		N				VIO200	바이오라이트	VIOlight	중국		지금 멤버십 가입하면 1개월 무료!		N									1	51	0								U	2022-11-07 10:10:48
...
```

3. 멀티 라인 파일 (multi line)
```
<<<begin>>>
<<<mapid>>>0A33908
<<<pname>>>ThinkPad 0A33908
<<<price>>>18900
<<<pgurl>>>https://www.lenovo.com/kr/ko/accessories-and-monitors/keyboards-and-mice/trackpoint-caps/CAP-Thinkpad-TrackPoint-Caps/p/0A33908?cid=kr:cse:DANAWA
<<<igurl>>>https://www.lenovo.com/medias/0A33908_200.png?context=bWFzdGVyfHJvb3R8MzcxMDR8aW1hZ2UvcG5nfGgxYi9oOWUvMTA5NDY0Nzc0NTc0MzgucG5nfDA4MzIxNmM0MWUwMDM3OGM0MmRiZjJjYjVkMGRmZTVkOTQyZDA3NGM3Y2Q2ODYxNzFmNjM3YWNiNjFlM2IyNmI
<<<cate1>>>SS&P
<<<cate2>>>
<<<cate3>>>TrackPoint
<<<cate4>>>ThinkPad
<<<caid1>>>
<<<model>>>
<<<maker>>>
<<<origi>>>
<<<deliv>>>0
<<<event>>>
<<<coupo>>>
<<<mpric>>>18900
<<<ftend>>>
...
```

## 파일 서처 서비스 API 가이드

`GET /search`

해당 업체 코드와 상품ID로 원본 파일의 내용을 조회합니다.

#### Request Parameter

|parameter|description|data type|example value|
|:---:|:---:|:---:|:---:|
|shopCode|업체 코드|string|TH201, ED302 , ..|
|productId|상품 아이디|string|2074042,7658077,2186943|
|startDate|날짜 From|string (YYYYMMDD)|20221110|
|endDate|날짜 To|string (YYYYMMDD)|20221111|
|renew|전체/갱신 여부|string|1 또는 2|

#### Return Value

- result : 검색 처리 결과 (success or fail)
- date : 검색 일자
- renew : 전체/갱신 여부
- data : 검색 내용
- timestamp : 해당 데이터 일자
- source : 원본 데이터 내용

GET /search?shopCode=EE301&productId=2073143136&startDate=20220326&&endDate=20220329

```JSON
[
    {
        "result": "success",
        "date": 20220926,
        "renew" : 1,
        "data" : [
            {
                "timestamp": 1664171131,
                "source" : "2073143136^이미용|화장품|메이크업(여성용)|볼터치/하이라이터^[현대백화점] [삼성카드7%할인~08/22]아워글래스 앰비언트 블러쉬 +무이자3개월^(주) 신세계인터네셔날^http://image.thehyundai.com/static/3/1/3/14/73/2073143136_0_600.jpg^http://www.thehyundai.com/front/pda/itemPtc.thd"
            },
            {
                "timestamp" : 1661929118,
                "source" : "2073143136^이미용|화장품|메이크업(여성용)|볼터치/하이라이터^[현대백화점] [삼성카드7%할인~08/22]아워글래스 앰비언트 블러쉬 +무이자3개월^(주) 신세계인터네셔날^http://image.thehyundai.com/static/3/1/3/14/73/2073143136_0_600.jpg^http://www.thehyundai.com/front/pda/itemPtc.thd"
            }
        ]
    },
    {
        "result" : "success",
        "date" : 20220925,
        "renew" : 1,
        "data" : [
            {
                "timestamp" : 1664171131,
                "source" :    "2073143136^이미용|화장품|메이크업(여성용)|볼터치/하이라이터^[현대백화점] [삼성카드7%할인~08/22]아워글래스 앰비언트 블러쉬 +무이자3개월^(주) 신세계인터네셔날^http://image.thehyundai.com/static/3/1/3/14/73/2073143136_0_600.jpg^http://www.thehyundai.com/front/pda/itemPtc.thd"
            },
            {
                "timestamp" : 1661929118,
                "source" :    "2073143136^이미용|화장품|메이크업(여성용)|볼터치/하이라이터^[현대백화점] [삼성카드7%할인~08/22]아워글래스 앰비언트 블러쉬 +무이자3개월^(주) 신세계인터네셔날^http://image.thehyundai.com/static/3/1/3/14/73/2073143136_0_600.jpg^http://www.thehyundai.com/front/pda/itemPtc.thd"
            }
        ]
    }
]
```

## 기술 블로그

https://danawalab.github.io/common/2022/11/10/big-file-searcher.html

